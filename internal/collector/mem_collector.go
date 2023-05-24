package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/mem"
)

type MemCollector struct {
	totalMemory     *prometheus.Desc
	usedMemory      *prometheus.Desc
	freeMemory      *prometheus.Desc
	cachedMemory    *prometheus.Desc
	bufferedMemory  *prometheus.Desc
	sharedMemory    *prometheus.Desc
	availableMemory *prometheus.Desc
}

func NewMemCollector() *MemCollector {
	return &MemCollector{
		totalMemory: prometheus.NewDesc("mem_total",
			"Total physical memory (RAM) in bytes.",
			nil, nil,
		),
		usedMemory: prometheus.NewDesc("mem_used",
			"Used memory in bytes.",
			nil, nil,
		),
		freeMemory: prometheus.NewDesc("mem_free",
			"Free memory in bytes.",
			nil, nil,
		),
		cachedMemory: prometheus.NewDesc("mem_cached",
			"Cached memory in bytes.",
			nil, nil,
		),
		bufferedMemory: prometheus.NewDesc("mem_buffered",
			"Buffered memory in bytes.",
			nil, nil,
		),
		sharedMemory: prometheus.NewDesc("mem_shared",
			"Shared memory in bytes.",
			nil, nil,
		),
		availableMemory: prometheus.NewDesc("mem_available",
			"Available memory in bytes.",
			nil, nil,
		),
	}
}

func (collector *MemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.totalMemory
	ch <- collector.usedMemory
	ch <- collector.freeMemory
	ch <- collector.cachedMemory
	ch <- collector.bufferedMemory
	ch <- collector.sharedMemory
	ch <- collector.availableMemory
}

func (collector *MemCollector) Collect(ch chan<- prometheus.Metric) {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.totalMemory, prometheus.GaugeValue, float64(vMem.Total))
	ch <- prometheus.MustNewConstMetric(collector.usedMemory, prometheus.GaugeValue, float64(vMem.Used))
	ch <- prometheus.MustNewConstMetric(collector.freeMemory, prometheus.GaugeValue, float64(vMem.Free))
	ch <- prometheus.MustNewConstMetric(collector.cachedMemory, prometheus.GaugeValue, float64(vMem.Cached))
	ch <- prometheus.MustNewConstMetric(collector.bufferedMemory, prometheus.GaugeValue, float64(vMem.Buffers))
	ch <- prometheus.MustNewConstMetric(collector.sharedMemory, prometheus.GaugeValue, float64(vMem.Shared))
	ch <- prometheus.MustNewConstMetric(collector.availableMemory, prometheus.GaugeValue, float64(vMem.Available))
}
