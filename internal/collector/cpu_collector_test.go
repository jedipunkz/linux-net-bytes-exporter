package collector

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func mockCPUPercent(interval time.Duration, percpu bool) ([]float64, error) {
	return []float64{50.0 * 100}, nil // Multiply by 100 to match the Collect function's behavior
}

func TestNewCPUCollector(t *testing.T) {
	collector := NewCPUCollector()
	if collector.cpuUsage == nil {
		t.Error("cpuUsage is nil")
	}
}

func TestDescribe(t *testing.T) {
	collector := NewCPUCollector()
	ch := make(chan *prometheus.Desc, 1)
	collector.Describe(ch)

	select {
	case desc := <-ch:
		if desc.String() != collector.cpuUsage.String() {
			t.Errorf("Expected description to be sent to the channel")
		}
	default:
		t.Errorf("Expected description to be sent to the channel")
	}
}

func findMetric(metricFamilies []*dto.MetricFamily, name string) *dto.Metric {
	for _, mf := range metricFamilies {
		if *mf.Name == name {
			for _, m := range mf.Metric {
				return m
			}
		}
	}
	return nil
}

func TestCollect(t *testing.T) {
	collector := NewCPUCollector()
	collector.CPUPercent = mockCPUPercent

	registerer := prometheus.NewPedanticRegistry()
	registerer.Register(collector)

	metricFamilies, err := registerer.Gather()
	if err != nil {
		t.Fatalf("Unexpected error gathering metrics: %v", err)
	}

	metric := findMetric(metricFamilies, "cpu_usage")
	if metric == nil {
		t.Fatal("cpu_usage metric not found")
	}

	value := metric.GetGauge().GetValue()
	if value != 50.0 {
		t.Errorf("Unexpected metric value: %v", value)
	}
}
