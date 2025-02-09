package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
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
	zpoolPoolSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_storage_size_bytes",
		},
		[]string{"pool"},
	)
	zpoolPoolUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_storage_usage_bytes",
		},
		[]string{"pool"},
	)
	zpoolPoolFree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_storage_free_bytes",
		},
		[]string{"pool"},
	)
	zpoolPoolFragmentation = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zpool_storage_fragmentation_percent",
		},
		[]string{"pool"},
	)
)

type metricGatherer func()
