package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func ParseTime(exp string) (time.Time, error) {
	re := regexp.MustCompile(`^(\d+)([dhm])$`)
	matches := re.FindStringSubmatch(exp)

	if len(matches) != 3 {
		return time.Time{}, fmt.Errorf("invalid expiration format: %s", exp)
	}

	value, err := strconv.Atoi(matches[1]) // Extract the numeric part, so 30d -> 30
	if err != nil {
		return time.Time{}, err
	}

	unit := matches[2] // Extract the unit (d, h, m), so 30d -> d

	var timeFromNow time.Time
	switch unit {
	case "d":
		timeFromNow = time.Now().Add(time.Duration(value) * 24 * time.Hour)
	case "h":
		timeFromNow = time.Now().Add(time.Duration(value) * time.Hour)
	case "m":
		timeFromNow = time.Now().Add(time.Duration(value) * time.Minute)
	default:
		return time.Time{}, fmt.Errorf("unsupported time unit: %s", unit)
	}
	return timeFromNow, nil
}
