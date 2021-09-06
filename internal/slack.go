package internal

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// GetTimeFromSlackTimeStamp returns time converted from string
func GetTimeFromSlackTimeStamp(ts string) (time.Time, error) {
	array := strings.Split(ts, ".")
	unixtime, err := strconv.ParseInt(array[0], 10, 64)
	if err != nil {
		return time.Time{}, errors.New("Failed to parse as int")
	}
	return time.Unix(unixtime, 0), nil
}
