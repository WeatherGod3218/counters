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
		return counters[i].LastReset.Timestamp.Before(counters[j].LastReset.Timestamp)
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

	reset := &database.Reset{
		Reporter:    "A Person",
		Instigator:  c.PostForm("reset-user"),
		Description: c.PostForm("reset-description"),
		Timestamp:   time.Now(),
	}

	newHistory := make([]database.Reset, 1)
	newHistory[0] = *reset

	counter := &database.Counter{
		Id:          bson.NewObjectID(),
		CreatedBy:   "A Person",
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		LastReset:   *reset,
		History:     newHistory,
	}

	counterId, err := database.CreateCounter(c, counter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/counters/"+counterId)
}
