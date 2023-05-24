package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUCollector struct {
	cpuUser   *prometheus.Desc
	cpuSystem *prometheus.Desc
	cpuIowait *prometheus.Desc
	cpuIrq    *prometheus.Desc
	cpuUsage  *prometheus.Desc
}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{
		cpuUsage: prometheus.NewDesc("cpu_usage",
			"CPU usage based on the idle time.",
			nil, nil,
		),
		cpuUser: prometheus.NewDesc(
			"cpu_user",
			"CPU utilization in the user mode as a percentage of total CPU capacity.",
			nil, nil,
		),
		cpuSystem: prometheus.NewDesc(
			"cpu_system",
			"CPU utilization in the system mode as a percentage of total CPU capacity.",
			nil, nil,
		),
		cpuIowait: prometheus.NewDesc(
			"cpu_iowait",
			"CPU time waiting for I/O operations to complete as a percentage of total CPU capacity.",
			nil, nil,
		),
		cpuIrq: prometheus.NewDesc(
			"cpu_irq",
			"CPU time servicing interrupts as a percentage of total CPU capacity.",
			nil, nil,
		),
	}
}

func (collector *CPUCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.cpuUsage
	ch <- collector.cpuUser
	ch <- collector.cpuSystem
	ch <- collector.cpuIowait
	ch <- collector.cpuIrq
}

func (collector *CPUCollector) Collect(ch chan<- prometheus.Metric) {
	cpuTimes, err := cpu.Times(false)
	if err != nil {
		return
	}

	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return
	}

	for _, cpuTime := range cpuTimes {
		ch <- prometheus.MustNewConstMetric(collector.cpuUsage, prometheus.GaugeValue, cpuPercent[0])
		ch <- prometheus.MustNewConstMetric(collector.cpuUser, prometheus.GaugeValue, cpuTime.User)
		ch <- prometheus.MustNewConstMetric(collector.cpuSystem, prometheus.GaugeValue, cpuTime.System)
		ch <- prometheus.MustNewConstMetric(collector.cpuIowait, prometheus.GaugeValue, cpuTime.Iowait)
		ch <- prometheus.MustNewConstMetric(collector.cpuIrq, prometheus.GaugeValue, cpuTime.Irq)
	}
}
