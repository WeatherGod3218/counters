package users

import (
	"errors"
	"os"
	"slices"
	"strconv"

	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/gin-gonic/gin"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
)

func IsEboard(user *csh_auth.UserInfo) bool {
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

func IsActiveRTP(user *csh_auth.UserInfo) bool {
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

func GetCSHAuth(c *gin.Context) (*csh_auth.UserInfo, error) {
	userAny, exists := c.Get("cshauth")

	if !exists {
		return nil, errors.New("unable to load csh auth")
	}

	userClaims, ok := userAny.(*csh_auth.Claims)
	if !ok {
		return nil, errors.New("unable to cast csh auth")
	}

	user := &csh_auth.UserInfo{
		Uuid:     userClaims.Uuid,
		Email:    userClaims.Email,
		Username: userClaims.Username,
		FullName: userClaims.FullName,
		Groups:   userClaims.Groups,
	}

	return user, nil
}
