package utils

import (
	"time"
)

const TimeFormatLayout = "2006-01-02T15:04:05Z"

func ParseTime(timeString string) (time.Time, error) {
	return time.Parse(TimeFormatLayout, timeString+"Z")
}
