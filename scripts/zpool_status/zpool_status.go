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

type ScanState struct {
	PercentDone float64
	ETA         time.Duration
}

type ZPool struct {
	Name  string
	Type  string
	State string
	Scan  ScanState
	Disks []ZDisk
}

var (
	zpoolPoolScan = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_scan_percent_done",
		},
		[]string{"pool"},
	)
	zpoolPoolScanEta = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_scan_eta",
		},
		[]string{"pool"},
	)
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

	s := new(ScanState)
	scanTmp := cmdOutput[strings.Index(cmdOutput, "scan:")+5 : strings.Index(cmdOutput, "config:")]
	if strings.Contains(scanTmp, "% done") {
		p := regexp.MustCompile("([0-9]{2}.[0-9]{2})% done").FindAllStringSubmatch(scanTmp, 1)[0][1]
		if s.PercentDone, err = strconv.ParseFloat(p, 64); err != nil {
			log.Fatal(err)
		}
	}
	if strings.Contains(scanTmp, "to go") {
		t := strings.Split(regexp.MustCompile("([0-9]{2}:[0-9]{2}:[0-9]{2}) to go").FindAllStringSubmatch(scanTmp, 1)[0][1], ":")
		s.ETA, _ = time.ParseDuration(fmt.Sprintf("%sh%sm%ss", t[0], t[1], t[2]))
	}
	out.Scan = *s

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

		zpoolPoolScan.WithLabelValues(data.Name).Set(data.Scan.PercentDone)
		zpoolPoolScanEta.WithLabelValues(data.Name).Set(float64(data.Scan.ETA.Milliseconds()))

	}
}

func init() {
	prometheus.Register(zpoolDiskState)
	prometheus.Register(zpoolDiskChecksumErrors)
	prometheus.Register(zpoolDiskReadErrors)
	prometheus.Register(zpoolDiskWriteErrors)
	prometheus.Register(zpoolPoolScan)
	prometheus.Register(zpoolPoolScanEta)
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
