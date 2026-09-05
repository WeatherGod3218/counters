package pages

import (
	"net/http"
	"sort"

	"github.com/WeatherGod3218/counters/database"
	"github.com/aws/smithy-go/logging"
	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetHomePage(c *gin.Context) {
	userAny, exists := c.Get("cshauth")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}
	user, ok := userAny.(*csh_auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}

	counters, err := database.GetAllCounters(c)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetHomePage"}).Warn("Unable to load counters!")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load counters!"})
		return
	}

	sort.Slice(counters, func(i, j int) bool {
		return counters[i].LastReset.Timestamp > counters[j].LastReset.Timestamp
	})

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Counters": counters,
		"Username": user.Username,
		"FullName": user.FullName,
		"EBoard":   IsEboard(user),
		"RTP":      IsActiveRTP(user),
	})
}
