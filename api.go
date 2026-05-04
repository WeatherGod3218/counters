package main

import (
	"net/http"

	"github.com/WeatherGod3218/counters/database"
	"github.com/WeatherGod3218/counters/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetHomePage(c *gin.Context) {
	// This is intentionally left unprotected
	// A user may be unable to vote but should still be able to see a list of polls

	counters, err := database.GetActiveCounters(c)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetHomePage"}).Fatal("error getting active counters")
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Counters": counters,
	})
}

func GetCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "create.tmpl", gin.H{
		"No": "Yes",
	})
}
