package users

import (
	"os"
	"slices"
	"strconv"

	"github.com/WeatherGod3218/counters/internal/logging"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
)

func IsEboard(user *csh_auth.Claims) bool {
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

func IsActiveRTP(user *csh_auth.Claims) bool {
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
