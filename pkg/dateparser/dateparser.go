package dateparser

import (
	"fmt"
	"time"
)

func Parse(data string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, data); err == nil {
		return t, err
	}
	if t, err := time.Parse(time.DateTime, data); err == nil {
		return t, err
	}
	if t, err := time.Parse(time.DateOnly, data); err == nil {
		return t, err
	}
	if t, err := time.Parse("2006-01-02T15:04:05", data); err == nil {
		return t, err
	}
	if t, err := time.Parse("2006-01-02T15:04:05-0700", data); err == nil {
		return t, err
	}
	if t, err := time.Parse("2006-01-02", data); err == nil {
		return t, err
	}
	if t, err := time.Parse("2006-01-02T15:04", data); err == nil {
		return t, err
	}
	if t, err := time.Parse("01.02.2006", data); err == nil {
		return t, err
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", data)
}
