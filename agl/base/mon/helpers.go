package mon

import "github.com/prometheus/client_golang/prometheus"

// Labels is an alias of prometheus.Labels.
type Labels = prometheus.Labels

// TrackInflight is a helper function to inc/dec a gauge.
func TrackInflight(g prometheus.Gauge) func() {
	g.Inc()
	did := false
	return func() {
		if did {
			// Some API allow the done to be run more than once.
			// But we don't want to decrement the gauge more than once.
			return
		}
		did = true
		g.Dec()
	}
}
