package internal

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// GetTimeFromSlackTimeStamp returns time converted from string
func GetTimeFromSlackTimeStamp(ts string) (time.Time, error) {
	array := strings.Split(ts, ".")
	unixtime, err := strconv.ParseInt(array[0], 10, 64)
	if err != nil {
		log.Println(err)
		return time.Time{}, errors.New("Failed to ParseInt()")
	}
	return time.Unix(unixtime, 0), nil
}
