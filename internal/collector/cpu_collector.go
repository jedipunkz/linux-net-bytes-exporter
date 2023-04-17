package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

type CPUCollector struct {
	cpuUsage   *prometheus.Desc
	CPUPercent func(interval time.Duration, percpu bool) ([]float64, error)
}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{
		cpuUsage: prometheus.NewDesc("cpu_usage",
			"CPU usage percentage.",
			nil, nil,
		),
		CPUPercent: cpu.Percent,
	}
}

func (collector *CPUCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.cpuUsage
}

func (collector *CPUCollector) Collect(ch chan<- prometheus.Metric) {
	cpuPercent, err := collector.CPUPercent(0, false)
	if err != nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(collector.cpuUsage, prometheus.GaugeValue, cpuPercent[0]/100)
}
