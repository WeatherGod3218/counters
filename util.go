package main

import (
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/WeatherGod3218/counters/logging"

	cshAuth "github.com/computersciencehouse/csh-auth/v2"
)

func IsEboard(user *cshAuth.Claims) bool {
	val, set := os.LookupEnv("DEV_FORCE_IS_EBOARD")
	if set {
		forced, err := strconv.ParseBool(val)
		if err == nil {
			logging.Logger.Info("Forced Override for Eboard")
			return forced
		}
		logging.Logger.Warn("FORCED EBOARD WAS MISTYPED, MAKE SURE IT'S EITHER true OR false")
	}

	return slices.Contains(user.Groups, "eboard")
}

func IsActiveRTP(user *cshAuth.Claims) bool {
	val, set := os.LookupEnv("DEV_FORCE_IS_RTP")
	if set {
		forced, err := strconv.ParseBool(val)
		if err == nil {
			logging.Logger.Info("Forced Override for RTP")
			return forced
		}
		logging.Logger.Warn("FORCED RTP WAS MISTYPED, MAKE SURE IT'S EITHER true OR false")
	}

	return slices.Contains(user.Groups, "active-rtp")
}

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
