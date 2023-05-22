package collector

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type CPUCollector struct {
	idle   *prometheus.Desc
	user   *prometheus.Desc
	system *prometheus.Desc
	iowait *prometheus.Desc
	irq    *prometheus.Desc
}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{
		idle: prometheus.NewDesc("cpu_usage_idle",
			"Percentage of CPU usage that is being used by idle process",
			nil, nil,
		),
		user: prometheus.NewDesc("cpu_usage_user",
			"Percentage of CPU usage that is being used by user level processes",
			nil, nil,
		),
		system: prometheus.NewDesc("cpu_usage_system",
			"Percentage of CPU usage that is being used by system/kernel level processes",
			nil, nil,
		),
		iowait: prometheus.NewDesc("cpu_usage_iowait",
			"Percentage of CPU usage that is being used by waiting for I/O operations to complete",
			nil, nil,
		),
		irq: prometheus.NewDesc("cpu_usage_irq",
			"Percentage of CPU usage that is being used by handling hardware interrupts",
			nil, nil,
		),
	}
}

func (collector *CPUCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.idle
	ch <- collector.user
	ch <- collector.system
	ch <- collector.iowait
	ch <- collector.irq
}

func (collector *CPUCollector) Collect(ch chan<- prometheus.Metric) {
	cpuSample, err := fetchCpuStats()
	if err != nil {
		return
	}

	fields := strings.Fields(cpuSample)
	if len(fields) < 9 {
		return
	}

	user, _ := time.ParseDuration(fields[1] + "ns")
	system, _ := time.ParseDuration(fields[3] + "ns")
	idle, _ := time.ParseDuration(fields[4] + "ns")
	iowait, _ := time.ParseDuration(fields[5] + "ns")
	irq, _ := time.ParseDuration(fields[6] + "ns")

	total := user + system + idle + iowait + irq

	ch <- prometheus.MustNewConstMetric(collector.idle, prometheus.GaugeValue, float64(idle)/float64(total))
	ch <- prometheus.MustNewConstMetric(collector.user, prometheus.GaugeValue, float64(user)/float64(total))
	ch <- prometheus.MustNewConstMetric(collector.system, prometheus.GaugeValue, float64(system)/float64(total))
	ch <- prometheus.MustNewConstMetric(collector.iowait, prometheus.GaugeValue, float64(iowait)/float64(total))
	ch <- prometheus.MustNewConstMetric(collector.irq, prometheus.GaugeValue, float64(irq)/float64(total))
}

func fetchCpuStats() (string, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 1024)
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(buf), nil
}
