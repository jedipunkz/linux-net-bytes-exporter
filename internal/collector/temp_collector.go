package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/host"
)

type TempCollector struct {
	temp *prometheus.Desc
}

func NewTempCollector() *TempCollector {
	return &TempCollector{
		temp: prometheus.NewDesc(
			"cpu_temperature_celsius",
			"Current temperature of the CPU.",
			[]string{"sensor"}, nil,
		),
	}
}

func (collector *TempCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.temp
}

func (collector *TempCollector) Collect(ch chan<- prometheus.Metric) {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return
	}

	for _, temp := range temps {
		ch <- prometheus.MustNewConstMetric(collector.temp, prometheus.GaugeValue, temp.Temperature, temp.SensorKey)
	}
}
