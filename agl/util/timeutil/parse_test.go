package timeutil_test

import (
	"testing"

	"github.com/leclerc04/go-tool/agl/util/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestTimeutil(t *testing.T) {
	runParseTimeMaybeWithTimezoneTest(t)
}

func runParseTimeMaybeWithTimezoneTest(t *testing.T) {
	timeStr := "12:00:01"
	parsedTime, err := timeutil.ParseTimeMaybeWithTimezone(timeStr)
	assert.NoError(t, err)
	assert.Equal(t, "0000-01-01 12:00:01 +0800 +0800", parsedTime.String())

	timeStr = "12:00:02+09:00"
	parsedTime, err = timeutil.ParseTimeMaybeWithTimezone(timeStr)
	assert.NoError(t, err)
	assert.Equal(t, "0000-01-01 12:00:02 +0900 +0900", parsedTime.String())

	timeStr = "10:25"
	_, err = timeutil.ParseTimeMaybeWithTimezone(timeStr)
	if err == nil {
		t.Error("expected error")
		return
	}
	assert.Contains(t, err.Error(), "InvalidArgument error: 10:25 is not a valid time string.")
}
