package internal

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"time"

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

var (
	zpoolStates = []ZpoolState{ZpoolStateOnline, ZpoolStateScrubbing, ZpoolStateDegraded, ZpoolStateDegradedScrubbing, ZpoolStateFaulted, ZpoolStateFaultedScrubbing}
	zpool_bin   = FindZpoolBin()
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

func statusMetrics() {
	var output ZpoolStatusOutput

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
				if disk.Path != "" && strings.Contains(disk.Path, "/dev/disk/by-") {
					tmp, err := locateDiskByUUID(disk.Path)
					if err != nil {
						log.Errorf("failed to locate disk via uuid: %v", err)
					} else {
						name = tmp
					}
				}
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
