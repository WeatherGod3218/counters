package pages

import (
	"net/http"
	"sort"

	"github.com/ComputerScienceHouse/counters/internal/database"
	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetHomePage(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
	}

	counters, err := database.GetCounters(c.Request.Context())

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetHomePage"}).Warn("Unable to load counters!")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load counters!"})
		return
	}

	sort.Slice(counters, func(i, j int) bool {
		return counters[i].ResetOccuredAt > counters[j].ResetOccuredAt
	})

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Counters": counters,
		"Username": user.Username,
		"FullName": user.FullName,
		"EBoard":   users.IsEboard(user),
		"RTP":      users.IsActiveRTP(user),
	})
}
