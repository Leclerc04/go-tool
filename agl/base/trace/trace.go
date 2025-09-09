package trace

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/golang/glog"
	"github.com/leclerc04/go-tool/agl/base/mon"
	"golang.org/x/net/trace"
)

var (
	metricStartedTotal = mon.NewCounterVec(
		"trace", "started_total", "A counter of number of traces started.",
		[]string{"family", "title"})
	metricInflightTraces = mon.NewGaugeVec(
		"trace", "inflight_traces", "Number of traces inflight.",
		[]string{"family", "title"})
	logAllTraceFlag = flag.Bool("agl_trace_to_stderr", false, "If true, print all traces to log, for diagnosis.")
)

var readLogAllTraceOnce sync.Once
var logAllTraceValue bool

func logAllTrace() bool {
	readLogAllTraceOnce.Do(func() {
		// 延迟执行，以支持外部dot文件等配置。
		v := os.Getenv("AGL_TRACE_TO_STDERR")
		if v == "" {
			logAllTraceValue = *logAllTraceFlag
			return
		}
		logAllTraceValue = (v == "1")
	})
	return logAllTraceValue
}

// T provides uniform interface to track short/long-lived object.
type T interface {
	Printf(format string, a ...interface{})
	SetError()
	Indent(delta int)
}

type traceWrap struct {
	t      trace.Trace
	indent string
}

func (t traceWrap) Printf(format string, a ...interface{}) {
	if logAllTrace() {
		fmt.Fprintln(
			os.Stderr, "[TRACE]",
			// 区分不同的trace span.
			fmt.Sprintf("%p", t.t),
			fmt.Sprintf(format, a...))
	}
	if t.indent != "" {
		t.t.LazyPrintf("%s%s", t.indent, fmt.Sprintf(format, a...))
		return
	}
	t.t.LazyPrintf(format, a...)
}

func (t traceWrap) SetError() {
	t.t.SetError()
}

func (t *traceWrap) Indent(delta int) {
	if delta > 0 {
		t.indent += "."
	} else {
		t.indent = t.indent[0 : len(t.indent)-1]
	}
}

var traceLabelCleanUpPattern = regexp.MustCompile(`[0-9]+\.`)

// New creates a new short-lived trace.
func New(family, title string) (T, func()) {
	titleForMetric := title
	if titleForMetric == "run" {
		titleForMetric = ""
	}
	titleForMetric = traceLabelCleanUpPattern.ReplaceAllString(titleForMetric, "")
	labels := mon.Labels{
		"family": family,
		"title":  titleForMetric,
	}
	metricStartedTotal.With(labels).Inc()
	dec := mon.TrackInflight(metricInflightTraces.With(labels))
	t := traceWrap{
		t:      trace.New(family, title),
		indent: "",
	}
	t.t.SetMaxEvents(2000)
	return &t, func() {
		t.t.Finish()
		dec()
	}
}

type noopTrace struct{}

func (t noopTrace) Printf(format string, a ...interface{}) {}
func (t noopTrace) SetError()                              {}
func (t noopTrace) Indent(delta int)                       {}

// Noop is a trace that does nothing.
var Noop T = &noopTrace{}

type glogWrap struct{}

func (t glogWrap) Printf(format string, a ...interface{}) {
	glog.InfoDepth(3, fmt.Sprintf(format, a...))
}
func (t glogWrap) SetError() {}

func (t glogWrap) Indent(delta int) {}

// GLogTrace should only be used for unit test.
var glogTrace T = &glogWrap{}
