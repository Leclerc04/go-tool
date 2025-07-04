package cron

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/leclecr04/go-tool/agl/base/sentry"
	"github.com/leclecr04/go-tool/agl/util/must"
	"github.com/robfig/cron"
)

// EveryMidnight is the crontab of every midnight. EveyMidnight(0) means mid night in Beijing timezone.
func EveryMidnight(offsetHour int) string {
	minute := rand.Intn(60)
	hour := (16 + offsetHour) % 24
	return fmt.Sprintf("0 %d %d * * *", minute, hour)
}

// ServeCrontab runs the given action according to the crontab schedule.
func ServeCrontab(ctx context.Context, runOnce bool, traceName string, crontab string, action func(ctx context.Context) error) {
	run := func(ctx context.Context) {
		runWithTraceAndSentry(ctx, "cron."+traceName, "run", !runOnce, action)
	}

	if runOnce {
		run(ctx)
		return
	}
	c := cron.New()
	ctx = sentry.CreateDetachedContext(ctx)
	must.Must(c.AddFunc(crontab, func() {
		run(ctx)
	}))
	c.Start()
}
