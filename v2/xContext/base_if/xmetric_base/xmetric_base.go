package xmetric_base

import "github.com/prometheus/client_golang/prometheus"

type MetricsIF interface {
	SummaryBy(name string, labels prometheus.Labels) prometheus.Observer
	CounterBy(name string, labels prometheus.Labels) prometheus.Counter
	HistogramBy(name string, labels prometheus.Labels) prometheus.Observer
	GaugeBy(name string, labels prometheus.Labels) prometheus.Gauge
	ErrorHandler(err error)
}
