package timeutil

import (
	"regexp"
	"time"

	"github.com/leclerc04/go-tool/agl/util/errs"
)

var timeWithoutZonePattern = regexp.MustCompile(`^([0-9]{2}:[0-9]{2}:[0-9]{2})$`)
var timeMaybeWithTimezone = regexp.MustCompile(`^([0-9]{2}:[0-9]{2}:[0-9]{2}[+-][0-9]{2}:[0-9]{2})|([0-9]{2}:[0-9]{2}:[0-9]{2})$`)

// ParseTimeMaybeWithTimezone try to parse time Str to time with timezone, if timeStr has no timezone, use +08:00
func ParseTimeMaybeWithTimezone(timeStr string) (t time.Time, err error) {
	if !timeMaybeWithTimezone.MatchString(timeStr) {
		return t, errs.InvalidArgument.Newf("%s is not a valid time string", timeStr)
	}
	if timeWithoutZonePattern.MatchString(timeStr) {
		timeStr = timeStr + "+08:00"
	}
	t, err = time.Parse("15:04:05Z07:00", timeStr)
	return t, err
}
