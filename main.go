package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jedipunkz/linux-tiny-exporter/internal/collector"
)

func main() {
	registry := prometheus.NewRegistry()

	netCollector := collector.NewNetCollector()
	registry.MustRegister(netCollector)

	cpuCollector := collector.NewCPUCollector()
	registry.MustRegister(cpuCollector)

	diskCollector := collector.NewDiskCollector()
	registry.MustRegister(diskCollector)

	memCollector := collector.NewMemCollector()
	registry.MustRegister(memCollector)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	fmt.Println("Starting exporter on :9101")
	if err := http.ListenAndServe(":9101", nil); err != nil {
		fmt.Printf("Error starting exporter: %v\n", err)
	}
}
