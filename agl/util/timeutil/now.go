package timeutil

import (
	"sync/atomic"
	"time"

	"github.com/leclecr04/go-tool/agl/util/must"
)

var (
	fakeTime atomic.Value
)

func init() {
	fakeTime.Store(time.Time{})
}

func SetFakeTime(t time.Time) {
	fakeTime.Store(t)
}

func AdvanceFakeTime(d time.Duration) {
	t := fakeTime.Load().(time.Time)
	t = t.Add(d)
	fakeTime.Store(t)
}

func SetAFakeTime() {
	t, e := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41.123456Z")
	must.Must(e)
	SetFakeTime(t)
}

func UnSetFakeTime() {
	fakeTime.Store(time.Time{})
}

func Now() time.Time {
	t := fakeTime.Load().(time.Time)
	if !t.IsZero() {
		return t
	}
	return time.Now().UTC()
}
