package utils

import (
	"time"
)

// ParseStrToDate converts string to time.Time with specific format
func ParseStrToDate(input string, format string) (time.Time, error) {
	datetime, err := time.Parse(format, input)
	if err != nil {
		return time.Time{}, err
	}
	return ToDate(datetime), nil
}

// ToDate converts datetime to date format, which sets hour, min, sec to zero
func ToDate(datetime time.Time) time.Time {
	return time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 0, 0, 0, 0, time.Local)
}

// GetDatetimeStart returns the earliest time of datetime
// ex: 2022-01-01 00:00:00
func GetDatetimeStart(datetime time.Time) time.Time {
	return ToDate(datetime)
}

// GetDatetimeStart returns the latest time of datetime
// ex: 2022-01-01 23:59:59
func GetDatetimeEnd(datetime time.Time) time.Time {
	return time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 23, 59, 59, 0, time.Local)
}
