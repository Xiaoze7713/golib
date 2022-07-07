package xmetric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	ServerName  string `json:"server_name" `
	SubSystem   string `json:"sub_system" `
	Environment string `json:"environment" `
	IDC         string `json:"idc" toml:"idc"`
	IP          string `json:"ip" `
}

type XMetric struct {
	mPool            *MetricPool
	opts             prometheus.Opts
	errorHandlerFunc func(err error)
}

var Handler *XMetric

func InitXMetric(conf *Config, errHdl func(err error)) {
	constLabels := prometheus.Labels{
		"env": conf.Environment,
		"idc": conf.IDC,
		"ip":  conf.IP,
	}
	Handler = &XMetric{
		errorHandlerFunc: errHdl,
		mPool:            NewMetricPool(conf.ServerName),
		opts: prometheus.Opts{
			Namespace:   conf.ServerName,
			Subsystem:   conf.SubSystem,
			ConstLabels: constLabels,
			Help:        "",
		},
	}
}

func (x *XMetric) ToLabelName(labels prometheus.Labels) []string {
	names := []string{}
	for k := range labels {
		names = append(names, k)
	}
	return names
}

func (x *XMetric) SummaryBy(name string, labels prometheus.Labels) prometheus.Observer {
	res := x.mPool.Summary(name, labels)
	if res == nil {
		x.mPool.NewMetrics(MetricSummary, x.opts, name, x.ToLabelName(labels))
	}
	res = x.mPool.Summary(name, labels)
	return res
}

func (x *XMetric) CounterBy(name string, labels prometheus.Labels) prometheus.Counter {
	res := x.mPool.Counter(name, labels)
	if res == nil {
		x.mPool.NewMetrics(MetricCounter, x.opts, name, x.ToLabelName(labels))
	}
	res = x.mPool.Counter(name, labels)
	return res
}

func (x *XMetric) HistogramBy(name string, labels prometheus.Labels) prometheus.Observer {
	res := x.mPool.Histogram(name, labels)
	if res == nil {
		x.mPool.NewMetrics(MetricHistogram, x.opts, name, x.ToLabelName(labels))
	}
	res = x.mPool.Histogram(name, labels)
	return res
}

func (x *XMetric) GaugeBy(name string, labels prometheus.Labels) prometheus.Gauge {
	res := x.mPool.Gauge(name, labels)
	if res == nil {
		x.mPool.NewMetrics(MetricGauge, x.opts, name, x.ToLabelName(labels))
	}
	res = x.mPool.Gauge(name, labels)
	return res
}

func (x *XMetric) ErrorHandler(err error) {
	if x.errorHandlerFunc != nil {
		x.errorHandlerFunc(err)
	}
}
