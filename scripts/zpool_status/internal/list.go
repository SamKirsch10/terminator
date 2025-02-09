package internal

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func zpoolList() (ZpoolListOutput, error) {
	var o ZpoolListOutput

	cmd := exec.Command(zpool_bin, "list", "-j")
	output, err := cmd.Output()
	if err != nil {
		log.Error("Failed to run `zpool list -j`")
		return o, err
	}

	err = json.Unmarshal(output, &o)

	return o, err
}

func listMetrics() {
	listing, err := zpoolList()
	if err != nil {
		log.Errorf("failed to get pool metrics: %v", err)
		return
	}
	for pool, cfg := range listing.Pools {
		size, err := parseDiskSize(cfg.Properties.Size.Value)
		if err != nil {
			log.Errorf("failed to get pool size: %v", err)
		} else {
			zpoolPoolSize.WithLabelValues(pool).Set(size)
		}

		size, err = parseDiskSize(cfg.Properties.Allocated.Value)
		if err != nil {
			log.Errorf("failed to get pool allocation: %v", err)
		} else {
			zpoolPoolUsage.WithLabelValues(pool).Set(size)
		}

		size, err = parseDiskSize(cfg.Properties.Free.Value)
		if err != nil {
			log.Errorf("failed to get pool free: %v", err)
		} else {
			zpoolPoolFree.WithLabelValues(pool).Set(size)
		}
		s := strings.Replace(cfg.Properties.Fragmentation.Value, "%", "", 1)
		p, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Errorf("failed to get pool fragmentation: %v", err)
		} else {
			zpoolPoolFragmentation.WithLabelValues(pool).Set(p)
		}

	}
}
