package util

import (
	"time"
)

func TranslateTime(inputTime string) int64 {
	timeZone, _ := time.LoadLocation("America/New_York")

	currentTime := time.Now()

	timeConverted, err := time.ParseInLocation("2006-01-02T15:04", inputTime, timeZone)
	if err != nil {
		timeConverted = currentTime
	}

	if timeConverted.After(currentTime) {
		timeConverted = currentTime
	}

	return timeConverted.Unix()
}
