package mon

import (
	"github.com/prometheus/client_golang/prometheus"
)

type optionValue struct {
	constLabels prometheus.Labels
}

func getOptions(options []Option) optionValue {
	opts := optionValue{}
	for _, v := range options {
		if v == nil {
			continue
		}
		v(&opts)
	}
	return opts
}

// Option specifies optional param to the metric function.
type Option func(*optionValue)

// ConstLabels sets constant labels to the metric.
func ConstLabels(labels prometheus.Labels) func(*optionValue) {
	return func(opts *optionValue) {
		opts.constLabels = labels
	}
}
