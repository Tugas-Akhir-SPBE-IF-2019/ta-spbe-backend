package prometheusinst

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Devices  prometheus.Gauge
	Upgrades *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}