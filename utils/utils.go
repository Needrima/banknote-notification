package utils

import (
	"time"
)

const TimeFormatLayout = "2006-01-02T15:04:05Z"

func ParseTimeStringToTime(timeString string) (time.Time, error) {
	return time.Parse(TimeFormatLayout, timeString+"Z")
}

func ParseTimeToString(t time.Time) string {
	timeStr := t.Format(time.RFC3339) // e.g "2022-06-21T11:43:24+01:06"
	return timeStr[:len(timeStr)-6]   // "2022-06-21T11:43:24"
}

func PeriodToScheduledTime(scheduledTime time.Time) float64 {
	elapsed := scheduledTime.Local().Sub(time.Now().Local())
	return elapsed.Seconds() - 3600
}
