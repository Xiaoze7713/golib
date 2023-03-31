package xmetric

import (
	"github.com/prometheus/client_golang/prometheus"
	"sort"
	"sync"
)

type MetricType int

const (
	MetricCounter   = MetricType(1)
	MetricGauge     = MetricType(2)
	MetricHistogram = MetricType(3)
	MetricSummary   = MetricType(4)
)

type MetricPool struct {
	name            string
	histogramVecSet map[string]*prometheus.HistogramVec
	summeryVecSet   map[string]*prometheus.SummaryVec
	counterVecSet   map[string]*prometheus.CounterVec
	gaugeVecSet     map[string]*prometheus.GaugeVec
	lock            *sync.RWMutex
}

func NewMetricPool(name string) *MetricPool {
	return &MetricPool{
		name:            name,
		histogramVecSet: map[string]*prometheus.HistogramVec{},
		summeryVecSet:   map[string]*prometheus.SummaryVec{},
		counterVecSet:   map[string]*prometheus.CounterVec{},
		gaugeVecSet:     map[string]*prometheus.GaugeVec{},
		lock:            &sync.RWMutex{},
	}
}

type Labels []string

func (m Labels) HashKey() string {
	var keys []string
	for _, labelName := range m {
		keys = append(keys, labelName)
	}
	sort.Strings(keys)
	hashStr := "hash_"
	for _, k := range keys {
		hashStr += "#" + k
	}
	return hashStr
}

func (u *MetricPool) Counter(name string, labels prometheus.Labels) prometheus.Counter {
	u.lock.RLock()
	defer u.lock.RUnlock()
	counterVec, ok := u.counterVecSet[name]
	if !ok {
		return nil
	}
	counter := counterVec.With(labels)
	return counter
}

func (u *MetricPool) Summary(name string, labels prometheus.Labels) prometheus.Observer {
	u.lock.RLock()
	defer u.lock.RUnlock()
	summaryVec, ok := u.summeryVecSet[name]
	if !ok {
		return nil
	}
	summary := summaryVec.With(labels)
	return summary
}

func (u *MetricPool) Histogram(name string, labels prometheus.Labels) prometheus.Observer {
	u.lock.RLock()
	defer u.lock.RUnlock()
	histogramVec, ok := u.histogramVecSet[name]
	if !ok {
		return nil
	}
	histogram := histogramVec.With(labels)
	return histogram
}

func (u *MetricPool) Gauge(name string, labels prometheus.Labels) prometheus.Gauge {
	u.lock.RLock()
	defer u.lock.RUnlock()
	gaugeVec, ok := u.gaugeVecSet[name]
	if !ok {
		return nil
	}
	gauge := gaugeVec.With(labels)
	return gauge
}

func (u *MetricPool) NewMetrics(metricType MetricType, opts prometheus.Opts, name string, labels Labels) bool {
	u.lock.Lock()
	defer u.lock.Unlock()
	switch metricType {
	case MetricCounter:
		u.counterVecSet[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      name,
		}, labels)
		prometheus.MustRegister(u.counterVecSet[name])
	case MetricGauge:
		u.gaugeVecSet[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      name,
		}, labels)
		prometheus.MustRegister(u.gaugeVecSet[name])
	case MetricHistogram:
		u.histogramVecSet[name] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      name,
			//Buckets:   buckets,
		}, labels)
		prometheus.MustRegister(u.histogramVecSet[name])
	case MetricSummary:
		u.summeryVecSet[name] = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      name,
			//Objectives: objectives,
		}, labels)
		prometheus.MustRegister(u.summeryVecSet[name])
	default:
		return false
	}
	return true
}

func (u *MetricPool) GetMetricsVectors() (collectors []prometheus.Collector) {
	collectors = []prometheus.Collector{}
	for _, metrics := range u.gaugeVecSet {
		collectors = append(collectors, metrics)
	}
	for _, metrics := range u.histogramVecSet {
		collectors = append(collectors, metrics)
	}
	for _, metrics := range u.counterVecSet {
		collectors = append(collectors, metrics)
	}
	for _, metrics := range u.summeryVecSet {
		collectors = append(collectors, metrics)
	}
	return collectors
}
