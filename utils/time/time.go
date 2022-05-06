package time

import "time"

// EmptyTime is a time whose unix timestamp is zero.
var EmptyTime = CreateEmptyTime()

// CreateEmptyTime creates an empty time whose unix timestamp is zero.
func CreateEmptyTime() time.Time {
	return time.Unix(0, 0)
}

// IsEmptyTime checks whether the time is empty, ie whose unix timestamp is zero.
func IsEmptyTime(t *time.Time) bool {
	return t.Unix() == 0
}
