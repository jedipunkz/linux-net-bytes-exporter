package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type netCollector struct {
	receivedBytes *prometheus.Desc
	transmitBytes *prometheus.Desc
}

func newNetCollector() *netCollector {
	return &netCollector{
		receivedBytes: prometheus.NewDesc("net_received_bytes_total",
			"Total number of received bytes by network interface.",
			[]string{"interface"}, nil,
		),
		transmitBytes: prometheus.NewDesc("net_transmit_bytes_total",
			"Total number of transmitted bytes by network interface.",
			[]string{"interface"}, nil,
		),
	}
}

func (collector *netCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.receivedBytes
	ch <- collector.transmitBytes
}

func (collector *netCollector) Collect(ch chan<- prometheus.Metric) {
	data, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if i < 2 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		interfaceName := strings.Trim(fields[0], ":")
		receivedBytes, _ := strconv.ParseFloat(fields[1], 64)
		transmitBytes, _ := strconv.ParseFloat(fields[9], 64)

		ch <- prometheus.MustNewConstMetric(collector.receivedBytes, prometheus.CounterValue, receivedBytes, interfaceName)
		ch <- prometheus.MustNewConstMetric(collector.transmitBytes, prometheus.CounterValue, transmitBytes, interfaceName)
	}
}

func main() {
	registry := prometheus.NewRegistry()
	nc := newNetCollector()
	registry.MustRegister(nc)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	fmt.Println("Starting exporter on :9100")
	if err := http.ListenAndServe(":9100", nil); err != nil {
		fmt.Printf("Error starting exporter: %v\n", err)
	}
}
