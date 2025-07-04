package buildinfo

import (
	"github.com/leclecr04/go-tool/agl/base/mon"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricBuildInfo = mon.NewGauge("base", "build_info",
		"the build info for a2 binary, one per processs.",
		mon.ConstLabels(prometheus.Labels{
			"build_time": BuildTime,
			"git_commit": GitCommit,
		}))
)

func init() {
	metricBuildInfo.Set(1)
}
