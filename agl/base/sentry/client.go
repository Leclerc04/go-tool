package sentry

import (
	raven "github.com/getsentry/raven-go"
	"github.com/leclerc04/go-tool/agl/base/trace"
	"github.com/leclerc04/go-tool/agl/util/buildinfo"
)

// NewClient creates a new sentry client.
func NewClient(dsn string) (*raven.Client, error) {
	client, err := raven.New(dsn)
	client.DropHandler = func(p *raven.Packet) {
		metricDroppedTotal.Inc()
		b, err := p.JSON()
		if err != nil {
			trace.Printf(droppedLog, "invalid packet: %v", err)
			return
		}
		trace.Printf(droppedLog, "%s", string(b))
	}
	client.SetRelease(buildinfo.Release())
	return client, err
}
