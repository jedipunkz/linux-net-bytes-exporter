package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/disk"
)

type DiskCollector struct {
	readBytes  *prometheus.Desc
	writeBytes *prometheus.Desc
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{
		readBytes: prometheus.NewDesc(
			"diskio_read_bytes",
			"Number of bytes read from disk",
			[]string{"disk"},
			nil,
		),
		writeBytes: prometheus.NewDesc(
			"diskio_write_bytes",
			"Number of bytes written to disk",
			[]string{"disk"},
			nil,
		),
	}
}

func (collector *DiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.readBytes
	ch <- collector.writeBytes
}

func (collector *DiskCollector) Collect(ch chan<- prometheus.Metric) {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return
	}

	for diskName, ioCounter := range ioCounters {
		ch <- prometheus.MustNewConstMetric(collector.readBytes, prometheus.GaugeValue, float64(ioCounter.ReadBytes), diskName)
		ch <- prometheus.MustNewConstMetric(collector.writeBytes, prometheus.GaugeValue, float64(ioCounter.WriteBytes), diskName)
	}
}
