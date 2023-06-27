package helper

import (
	"time"
)

// RFC3339 format layout for time.
const TimeFormatLayout = "2006-01-02T15:04:05Z"

// ParseTimeStringToTime converts a timestring to time in RFC3339 format.
func ParseTimeStringToTime(timeString string) (time.Time, error) {
	return time.Parse(TimeFormatLayout, timeString+"Z") // trailing "Z" allows for parsing the timeString
}

// ParseTimeToString converts a time to string format.
func ParseTimeToString(t time.Time) string {
	timeStr := t.Format(time.RFC3339) // "2022-06-21T11:43:24+01:06"
	return timeStr[:len(timeStr)-6]   // "2022-06-21T11:43:24"
}

// PeriodToScheduledTime calculates the period between the current local time and scheduledTime and
// returns the value in seconds. It subtracts an hour (3600 seconds) from the elapsed time before returning the result
// because the elapsed time is an hour ahead of local time.
func PeriodToScheduledTime(scheduledTime string) float64 {
	t, _ := ParseTimeStringToTime(scheduledTime)
	elapsed := t.Local().Sub(time.Now().Local())
	return elapsed.Seconds() - 3600
}
