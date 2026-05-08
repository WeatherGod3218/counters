package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/WeatherGod3218/counters/database"
	"github.com/WeatherGod3218/counters/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetHomePage(c *gin.Context) {
	// This is intentionally left unprotected
	// A user may be unable to vote but should still be able to see a list of polls

	counters, err := database.GetAllCounters(c)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetHomePage"}).Fatal("error getting active counters")
	}

	sort.Slice(counters, func(i, j int) bool {
		return counters[i].LastReset.Timestamp.After(counters[j].LastReset.Timestamp)
	})

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Counters": counters,
	})
}

func GetCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "create.tmpl", gin.H{
		"No": "Yes",
	})
}

func GetCounterId(c *gin.Context) {
	counter, err := database.GetCounterFromId(c, c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.HTML(200, "counter.tmpl", gin.H{
		"Id":          counter.Id.Hex(),
		"Title":       counter.Title,
		"Description": counter.Description,
	})
}

func CreateCounter(c *gin.Context) {
	// convert time to New York
	logging.Logger.WithFields(logrus.Fields{"timeOnForm": c.PostForm("reset-time"), "module": "api", "method": "GetHomePage"}).Info("Time Registered On Form")

	timeZone, _ := time.LoadLocation("America/New_York")
	resetTime := c.PostForm("reset-time")

	currentTime := time.Now()

	timeConverted, err := time.ParseInLocation("2006-01-02T15:04", resetTime, timeZone)
	if err != nil {
		timeConverted = currentTime
	}

	if timeConverted.After(currentTime) {
		timeConverted = currentTime
	}

	reset := &database.Reset{
		Reporter:    "A Person",
		Instigator:  c.PostForm("reset-user"),
		Description: c.PostForm("reset-description"),
		Timestamp:   timeConverted,
	}

	newHistory := make([]database.Reset, 1)
	newHistory[0] = *reset

	newId := bson.NewObjectID()

	counter := &database.Counter{
		Id:          newId,
		CreatedBy:   "A Person",
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		LastReset:   *reset,
	}

	history := &database.History{
		Id:      newId,
		History: newHistory,
	}

	counterError := database.CreateCounter(c, counter)
	if counterError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": counterError.Error()})
		return
	}

	historyError := database.CreateHistory(c, history)
	if historyError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": historyError.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/counters/"+newId.Hex())
}
