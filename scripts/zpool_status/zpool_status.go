///usr/bin/env true; exec /usr/bin/env go run "$0" "$@"

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hairyhenderson/go-which"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	gatherSleep = 15 * time.Second
)

type VDevType string
type ZpoolState string

const (
	VDevTypeDisk                = VDevType("disk")
	VDevTypeRaidz               = VDevType("raidz")
	VDevTypeMirror              = VDevType("mirror")
	ZpoolStateOnline            = ZpoolState("ONLINE")
	ZpoolStateScrubbing         = ZpoolState("ONLINE/SCRUBBING")
	ZpoolStateDegraded          = ZpoolState("DEGRADED")
	ZpoolStateDegradedScrubbing = ZpoolState("DEGRADED/SCRUBBING")
	ZpoolStateFaulted           = ZpoolState("FAULTED")
	ZpoolStateFaultedScrubbing  = ZpoolState("FAULTED/SCRUBBING")
)

type Disk struct {
	Name           string   `json:"name"`
	VdevType       VDevType `json:"vdev_type"`
	GUID           string   `json:"guid"`
	Class          string   `json:"class"`
	State          string   `json:"state"`
	AllocSpace     string   `json:"alloc_space"`
	TotalSpace     string   `json:"total_space"`
	DefSpace       string   `json:"def_space"`
	ReadErrors     string   `json:"read_errors"`
	WriteErrors    string   `json:"write_errors"`
	ChecksumErrors string   `json:"checksum_errors"`
}

type VDev struct {
	Disk
	DiskMakeup map[string]Disk `json:"vdevs"`
}

type ZPool struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	PoolGUID   string `json:"pool_guid"`
	Txg        string `json:"txg"`
	SpaVersion string `json:"spa_version"`
	ZplVersion string `json:"zpl_version"`
	ErrorCount string `json:"error_count"`
	ScanStats  struct {
		Function           string `json:"function"`
		State              string `json:"state"`
		StartTime          string `json:"start_time"`
		EndTime            string `json:"end_time"`
		ToExamine          string `json:"to_examine"`
		Examined           string `json:"examined"`
		Skipped            string `json:"skipped"`
		Processed          string `json:"processed"`
		Errors             string `json:"errors"`
		BytesPerScan       string `json:"bytes_per_scan"`
		PassStart          string `json:"pass_start"`
		ScrubPause         string `json:"scrub_pause"`
		ScrubSpentPaused   string `json:"scrub_spent_paused"`
		IssuedBytesPerScan string `json:"issued_bytes_per_scan"`
		Issued             string `json:"issued"`
	} `json:"scan_stats"`
	Vdevs map[string]VDev `json:"vdevs"`
}

type ZpoolOutput struct {
	OutputVersion struct {
		Command   string `json:"command"`
		VersMajor int    `json:"vers_major"`
		VersMinor int    `json:"vers_minor"`
	} `json:"output_version"`
	Pools map[string]ZPool `json:"pools"`
}

var (
	zpool_bin = findZpoolBin()

	zpoolStates = []ZpoolState{ZpoolStateOnline, ZpoolStateScrubbing, ZpoolStateDegraded, ZpoolStateDegradedScrubbing, ZpoolStateFaulted, ZpoolStateFaultedScrubbing}

	zpoolPoolScan = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_scan_percent_done",
		},
		[]string{"pool"},
	)
	zpoolPoolState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_state",
		},
		[]string{"pool", "state"},
	)
	zpoolDiskState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_disk_state",
		},
		[]string{"pool", "vdev", "disk", "state"},
	)
	zpoolDiskReadErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_disk_read_errors",
		},
		[]string{"pool", "vdev", "disk"},
	)
	zpoolDiskWriteErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_disk_write_errors",
		},
		[]string{"pool", "vdev", "disk"},
	)
	zpoolDiskChecksumErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_disk_chksum_errors",
		},
		[]string{"pool", "vdev", "disk"},
	)
)

func zpoolStatus() (string, error) {
	cmd := exec.Command(zpool_bin, "status", "-j", "--json-flat-vdevs")
	output, err := cmd.Output()
	if err != nil {
		log.Error("Failed to run `zpool status tank -j --json-flat-vdevs`")
		return "", err
	}
	return string(output), nil
}

func parseDiskSize(sizeStr string) (float64, error) {
	if sizeStr == "" {
		return 0, errors.New("no string size detected!")
	}
	sizeStr = strings.TrimSpace(sizeStr)
	log.Debug("trying to parse disk size: ", sizeStr)

	var unit float64 = 1
	switch strings.ToUpper(string(sizeStr[len(sizeStr)-1])) {
	case "K":
		unit = 1024
	case "M":
		unit = 1024 * 1024
	case "G":
		unit = 1024 * 1024 * 1024
	case "T":
		unit = 1024 * 1024 * 1024 * 1024
	}

	sizeStr = sizeStr[:len(sizeStr)-1]

	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, err
	}

	return size * unit, nil
}

func dirtyStringToFloat(num string) float64 {
	o, _ := strconv.ParseFloat(num, 64)
	return o
}

func parseScrubPercent(data ZPool) (float64, error) {
	total, err := parseDiskSize(data.ScanStats.ToExamine)
	if err != nil {
		log.Error("failed to calculate bytes total size for percent done metric")
		return 0, err
	}
	done, err := parseDiskSize(data.ScanStats.Issued)
	if err != nil {
		log.Error("failed to calculate bytes done size for percent done metric")
		return 0, err
	}
	log.Debugf("parcing percent with values %f / %f", done, total)

	return (done / total) * 100, nil
}

func getPoolState(data ZPool) ZpoolState {
	isScrubbing := data.ScanStats.Function == "SCRUB" && data.ScanStats.State != "FINISHED"
	state := ZpoolState(data.State)

	if !isScrubbing {
		return state
	}
	switch state {
	case ZpoolStateOnline:
		return ZpoolStateScrubbing
	case ZpoolStateDegraded:
		return ZpoolStateDegradedScrubbing
	case ZpoolStateFaulted:
		return ZpoolStateFaultedScrubbing
	default:
		panic("unknown state during scrubbing!")
	}
}

func gatherMetrics() {
	log.Info("Gathering Metrics...")
	var output ZpoolOutput

	if cmdOutput, err := zpoolStatus(); err != nil {
		log.Fatal(err)
	} else {
		if err = json.Unmarshal([]byte(cmdOutput), &output); err != nil {
			log.Fatal(err)
		}
	}

	for pool, data := range output.Pools {
		log.Debugf("Checking pool %s", pool)
		for _, state := range zpoolStates {
			var status float64
			state = getPoolState(data)
			if ZpoolState(data.State) == state {
				status = 1
			}
			zpoolPoolState.WithLabelValues(data.Name, string(state)).Set(status)
		}
		for vDevName, vdev := range data.Vdevs {
			if vDevName == pool {
				continue
			}
			for name, disk := range vdev.DiskMakeup {
				for _, state := range zpoolStates {
					var status float64
					if ZpoolState(data.State) == state {
						status = 1
					}
					zpoolDiskState.WithLabelValues(data.Name, vDevName, name, string(state)).Set(status)
				}
				zpoolDiskChecksumErrors.WithLabelValues(data.Name, vDevName, name).Set(dirtyStringToFloat(disk.ChecksumErrors))
				zpoolDiskReadErrors.WithLabelValues(data.Name, vDevName, name).Set(dirtyStringToFloat(disk.ReadErrors))
				zpoolDiskWriteErrors.WithLabelValues(data.Name, vDevName, name).Set(dirtyStringToFloat(disk.WriteErrors))
			}
		}

		log.Debugf("Trying to parse scrub status: %s issued / %s toExamine", data.ScanStats.Issued, data.ScanStats.ToExamine)
		percent, err := parseScrubPercent(data)
		if err != nil {
			log.Errorf("failed to calculate scrub percent: %v", err)
		} else {
			zpoolPoolScan.WithLabelValues(data.Name).Set(percent)
		}

	}
}

func init() {
	prometheus.Register(zpoolDiskState)
	prometheus.Register(zpoolDiskChecksumErrors)
	prometheus.Register(zpoolDiskReadErrors)
	prometheus.Register(zpoolDiskWriteErrors)
	prometheus.Register(zpoolPoolScan)
	prometheus.Register(zpoolPoolState)
}

func findZpoolBin() string {
	return which.Which("zpool")
}

func main() {
	port := flag.String("port", "9000", "Port to listen on")
	lvl := flag.String("log-lvl", "INFO", "Log level")
	flag.Parse()

	loglvl := log.WarnLevel
	switch strings.ToUpper(*lvl) {
	case "INFO":
		loglvl = log.InfoLevel
	case "DEBUG":
		loglvl = log.DebugLevel
	case "WARN":
		loglvl = log.WarnLevel
	case "ERROR":
		loglvl = log.ErrorLevel
	default:
		panic("unknown log level. try `INFO`, `DEBUG`, `WARN`, or `ERROR`")
	}
	log.SetLevel(loglvl)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	ctx, cancel := context.WithCancel(context.Background())

	if zpool_bin == "" {
		log.Fatal("failed to find `zpool` binary. Exiting")
	}

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

	log.Info("Serving requests on port " + *port)
	err := http.ListenAndServe(":"+*port, router)
	log.Fatal(err)

}
