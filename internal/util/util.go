package util

import (
	"time"
)

func TranslateTime(inputTime string) int64 {
	timeZone, _ := time.LoadLocation("America/New_York")

	var finalTime int64 = 0
	currentTime := time.Now().Unix()

	timeConverted, err := time.ParseInLocation("2006-01-02T15:04", inputTime, timeZone)

	finalTime = timeConverted.Unix()

	if err != nil {
		finalTime = currentTime
	}

	if finalTime > currentTime {
		finalTime = currentTime
	}

	return finalTime
}
