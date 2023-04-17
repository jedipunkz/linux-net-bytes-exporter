package collector

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type NetCollector struct {
	receivedBytesDiff   *prometheus.Desc
	transmitBytesDiff   *prometheus.Desc
	receivedPacketsDiff *prometheus.Desc
	transmitPacketsDiff *prometheus.Desc
	lastReceivedBytes   map[string]float64
	lastTransmitBytes   map[string]float64
	lastReceivedPackets map[string]float64
	lastTransmitPackets map[string]float64
}

func NewNetCollector() *NetCollector {
	return &NetCollector{
		receivedBytesDiff: prometheus.NewDesc("net_received_bytes_diff",
			"Received bytes by network interface since the last scrape.",
			[]string{"interface"}, nil,
		),
		transmitBytesDiff: prometheus.NewDesc("net_transmit_bytes_diff",
			"Transmitted bytes by network interface since the last scrape.",
			[]string{"interface"}, nil,
		),
		receivedPacketsDiff: prometheus.NewDesc("net_received_packets_diff",
			"Received packets by network interface since the last scrape.",
			[]string{"interface"}, nil,
		),
		transmitPacketsDiff: prometheus.NewDesc("net_transmit_packets_diff",
			"Transmitted packets by network interface since the last scrape.",
			[]string{"interface"}, nil,
		),
		lastReceivedBytes:   make(map[string]float64),
		lastTransmitBytes:   make(map[string]float64),
		lastReceivedPackets: make(map[string]float64),
		lastTransmitPackets: make(map[string]float64),
	}
}

func (collector *NetCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.receivedBytesDiff
	ch <- collector.transmitBytesDiff
	ch <- collector.receivedPacketsDiff
	ch <- collector.transmitPacketsDiff
}

func (collector *NetCollector) Collect(ch chan<- prometheus.Metric) {
	data, err := getNetDevice()
	if err != nil {
		return
	}

	lines := strings.Split(data, "\n")

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
		receivedPackets, _ := strconv.ParseFloat(fields[2], 64)
		transmitPackets, _ := strconv.ParseFloat(fields[10], 64)

		if lastReceivedBytes, ok := collector.lastReceivedBytes[interfaceName]; ok {
			receivedBytesDiff := receivedBytes - lastReceivedBytes
			ch <- prometheus.MustNewConstMetric(collector.receivedBytesDiff, prometheus.GaugeValue, receivedBytesDiff, interfaceName)
		}
		if lastTransmitBytes, ok := collector.lastTransmitBytes[interfaceName]; ok {
			transmitBytesDiff := transmitBytes - lastTransmitBytes
			ch <- prometheus.MustNewConstMetric(collector.transmitBytesDiff, prometheus.GaugeValue, transmitBytesDiff, interfaceName)
		}
		if lastReceivedPackets, ok := collector.lastReceivedPackets[interfaceName]; ok {
			receivedPacketsDiff := receivedPackets - lastReceivedPackets
			ch <- prometheus.MustNewConstMetric(collector.receivedPacketsDiff, prometheus.GaugeValue, receivedPacketsDiff, interfaceName)
		}
		if lastTransmitPackets, ok := collector.lastTransmitPackets[interfaceName]; ok {
			transmitPacketsDiff := transmitPackets - lastTransmitPackets
			ch <- prometheus.MustNewConstMetric(collector.transmitPacketsDiff, prometheus.GaugeValue, transmitPacketsDiff, interfaceName)
		}

		collector.lastReceivedBytes[interfaceName] = receivedBytes
		collector.lastTransmitBytes[interfaceName] = transmitBytes
		collector.lastReceivedPackets[interfaceName] = receivedPackets
		collector.lastTransmitPackets[interfaceName] = transmitPackets
	}
}

func getNetDevice() (string, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, file)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
