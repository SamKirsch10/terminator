///usr/bin/env true; exec /usr/bin/env go run "$0" "$@"

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	gatherSleep = 15 * time.Second
)

type ZDisk struct {
	Disk        string
	State       string
	ReadErrors  int
	WriteErrors int
	CKSumErrors int
}

type ZPool struct {
	Name  string
	Type  string
	State string
	Disks []ZDisk
}

var (
	zpoolDiskState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_disk_state",
		},
		[]string{"pool", "disk", "state"},
	)
	zpoolDiskReadErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_disk_read_errors",
		},
		[]string{"pool", "disk"},
	)
	zpoolDiskWriteErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_disk_write_errors",
		},
		[]string{"pool", "disk"},
	)
	zpoolDiskChecksumErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_disk_chksum_errors",
		},
		[]string{"pool", "disk"},
	)
	zpoolStates = []string{"ONLINE", "DEGRADED", "FAULTED"}
)

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, strings.Replace(str, " ", "", -1))
		}
	}
	return r
}

func gatherData() *ZPool {
	out := new(ZPool)

	cmd := exec.Command("/usr/sbin/zpool", "status", "tank")
	output, err := cmd.Output()
	if err != nil {
		log.Println("[ERROR] Failed to run `zpool status tank`")
		return out
	}
	cmdOutput := string(output)

	out.Name = regexp.MustCompile("pool:\\s(.*)\n").FindAllStringSubmatch(cmdOutput, 1)[0][1]
	out.State = regexp.MustCompile("state:\\s(.*)\n").FindAllStringSubmatch(cmdOutput, 1)[0][1]

	tmp := cmdOutput[strings.Index(cmdOutput, "config:")+7 : strings.Index(cmdOutput, "errors:")]

	var dataFound bool
	var poolFound bool
	for _, line := range strings.Split(tmp, "\n") {
		d := new(ZDisk)
		if line == "" {
			continue
		}
		elements := delete_empty(strings.Split(strings.Replace(line, "\t", "", -1), "  "))
		if elements[0] == out.Name {
			poolFound = true
		} else if elements[0][0:4] == "raid" {
			dataFound = true
			continue
		}

		if dataFound && poolFound {
			d.Disk = elements[0]
			d.State = elements[1]
			d.ReadErrors, _ = strconv.Atoi(elements[2])
			d.WriteErrors, _ = strconv.Atoi(elements[3])
			d.CKSumErrors, _ = strconv.Atoi(elements[4])
			out.Disks = append(out.Disks, *d)
		}
	}

	return out
}

func gatherMetrics() {
	log.Println("[INFO] Gathering Metrics...")
	data := gatherData()
	for _, disk := range data.Disks {
		for _, state := range zpoolStates {
			if disk.State == state {
				zpoolDiskState.WithLabelValues(data.Name, disk.Disk, state).Set(1)
			} else {
				zpoolDiskState.WithLabelValues(data.Name, disk.Disk, state).Set(0)
			}
		}
		zpoolDiskChecksumErrors.WithLabelValues(data.Name, disk.Disk).Set(float64(disk.CKSumErrors))
		zpoolDiskReadErrors.WithLabelValues(data.Name, disk.Disk).Set(float64(disk.ReadErrors))
		zpoolDiskWriteErrors.WithLabelValues(data.Name, disk.Disk).Set(float64(disk.WriteErrors))

	}
}

func init() {
	prometheus.Register(zpoolDiskState)
	prometheus.Register(zpoolDiskChecksumErrors)
	prometheus.Register(zpoolDiskReadErrors)
	prometheus.Register(zpoolDiskWriteErrors)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		gatherMetrics()

		ticker := time.NewTicker(gatherSleep)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				gatherMetrics()
			}
		}
	}(ctx)
	defer cancel()

	router := mux.NewRouter()

	// Prometheus endpoint
	router.Path("/metrics").Handler(promhttp.Handler())

	fmt.Println("Serving requests on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)

}
