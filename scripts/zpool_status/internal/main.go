package internal

import (
	"context"
	"sync"
	"time"

	"github.com/hairyhenderson/go-which"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.Register(zpoolDiskState)
	prometheus.Register(zpoolDiskChecksumErrors)
	prometheus.Register(zpoolDiskReadErrors)
	prometheus.Register(zpoolDiskWriteErrors)
	prometheus.Register(zpoolPoolScan)
	prometheus.Register(zpoolPoolState)
}

func FindZpoolBin() string {
	return which.Which("zpool")
}

func Run(ctx context.Context) {

	gather()

	ticker := time.NewTicker(gatherSleep)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			gather()
		}
	}
}

func gather() {
	wg := sync.WaitGroup{}
	for _, f := range []metricGatherer{statusMetrics, listMetrics} {
		f := f
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
}
