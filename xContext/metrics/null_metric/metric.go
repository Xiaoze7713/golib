package null_metric

import "github.com/prometheus/client_golang/prometheus"

type MetricsNull struct{}

var ops = prometheus.Opts{
	Namespace:   "default",
	Subsystem:   "none",
	Name:        "null",
	Help:        "null",
	ConstLabels: nil,
}

func (MetricsNull) SummaryBy(name string, labels prometheus.Labels) prometheus.Observer {
	return prometheus.NewSummary(prometheus.SummaryOpts{Name: "summary", Namespace: ops.Namespace})
}
func (MetricsNull) CounterBy(name string, labels prometheus.Labels) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{Name: "cnt", Namespace: ops.Namespace})
}
func (MetricsNull) HistogramBy(name string, labels prometheus.Labels) prometheus.Observer {
	return prometheus.NewSummary(prometheus.SummaryOpts{Name: "his", Namespace: ops.Namespace})
}
func (MetricsNull) GaugeBy(name string, labels prometheus.Labels) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{Name: "gauge", Namespace: ops.Namespace})
}
func (MetricsNull) ErrorHandler(err error) {}
