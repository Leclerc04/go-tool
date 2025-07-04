package mon

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Name retruns a metric name.
func Name(module, name string) string {
	return prometheus.BuildFQName("a2", module, name)
}

// NewGauge create and register prometheus.Gauge.
func NewGauge(module, name, help string, options ...Option) prometheus.Gauge {
	opts := getOptions(options)
	v := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        Name(module, name),
		Help:        help,
		ConstLabels: opts.constLabels,
	})
	prometheus.MustRegister(v)
	return v
}

// NewGaugeVec create and register prometheus.GaugeVec.
func NewGaugeVec(module, name, help string, labels []string, options ...Option) *prometheus.GaugeVec {
	opts := getOptions(options)
	v := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        Name(module, name),
		Help:        help,
		ConstLabels: opts.constLabels,
	}, labels)
	prometheus.MustRegister(v)
	return v
}

// NewCounter create and register prometheus.Counter.
func NewCounter(module, name, help string, options ...Option) prometheus.Counter {
	opts := getOptions(options)
	v := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        Name(module, name),
		Help:        help,
		ConstLabels: opts.constLabels,
	})
	prometheus.MustRegister(v)
	return v
}

// NewCounterVec creates and registers prometheus.CounterVec.
func NewCounterVec(module, name, help string, labels []string, options ...Option) *prometheus.CounterVec {
	opts := getOptions(options)
	v := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        Name(module, name),
		Help:        help,
		ConstLabels: opts.constLabels,
	}, labels)
	prometheus.MustRegister(v)
	return v
}

// NewHistogramVec creates and registers prometheus.HistogramVec.
func NewHistogramVec(module, name, help string, buckets []float64, labels []string, options ...Option) *prometheus.HistogramVec {
	opts := getOptions(options)
	v := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        Name(module, name),
		Help:        help,
		Buckets:     buckets,
		ConstLabels: opts.constLabels,
	}, labels)
	prometheus.MustRegister(v)
	return v
}
