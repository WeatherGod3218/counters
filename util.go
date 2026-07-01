package main

import (
	"slices"
	"time"

	cshAuth "github.com/computersciencehouse/csh-auth/v2"
)

func IsEboard(user *cshAuth.Claims) bool {
	return slices.Contains(user.Groups, "eboard")
}

func IsActiveRTP(user *cshAuth.Claims) bool {
	return slices.Contains(user.Groups, "active-rtp")
}

func TranslateTime(inputTime string) time.Time {
	timeZone, _ := time.LoadLocation("America/New_York")

	currentTime := time.Now()

	timeConverted, err := time.ParseInLocation("2006-01-02T15:04", inputTime, timeZone)
	if err != nil {
		timeConverted = currentTime
	}

	if timeConverted.After(currentTime) {
		timeConverted = currentTime
	}
	return timeConverted
}
