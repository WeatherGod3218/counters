package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/WeatherGod3218/counters/database"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const tempName string = "A Person!"

func translateTime(inputTime string) time.Time {
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

func GetHomePage(c *gin.Context) {
	// This is intentionally left unprotected
	// A user may be unable to vote but should still be able to see a list of polls

	counters, err := database.GetAllCounters(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sort.Slice(counters, func(i, j int) bool {
		return counters[i].LastReset.Timestamp.After(counters[j].LastReset.Timestamp)
	})

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Counters": counters,
	})
}

func GetCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "create.tmpl", gin.H{})
}

func GetResetPage(c *gin.Context) {
	idToUse := c.Param("id")
	c.HTML(http.StatusOK, "reset.tmpl", gin.H{
		"Id": idToUse,
	})
}

func GetCounterId(c *gin.Context) {
	idToUse := c.Param("id")

	counter, cError := database.GetCounterFromId(c, idToUse)
	if cError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": cError.Error()})
		return
	}

	history, hError := database.GetHistoryFromId(c, idToUse)
	if hError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": hError.Error()})
		return
	}

	c.HTML(200, "counter.tmpl", gin.H{
		"Id":          counter.Id.Hex(),
		"Title":       counter.Title,
		"Description": counter.Description,
		"Timestamp":   counter.LastReset.Timestamp,
		"History":     history.History,
	})
}

func ResetCounter(c *gin.Context) {
	resetTime := translateTime(c.PostForm("reset-time"))

	reset := &database.Reset{
		Reporter:    tempName,
		Instigator:  c.PostForm("reset-user"),
		Description: c.PostForm("reset-description"),
		Timestamp:   resetTime,
	}

	counter, err := database.GetCounterFromId(c, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	counter.Reset(c, reset)

	c.Redirect(http.StatusFound, "/counters/"+counter.IDHex())
}

func CreateCounter(c *gin.Context) {
	resetTime := translateTime(c.PostForm("reset-time"))

	reset := &database.Reset{
		Reporter:    tempName,
		Instigator:  c.PostForm("reset-user"),
		Description: c.PostForm("reset-description"),
		Timestamp:   resetTime,
	}

	counter := &database.Counter{
		Id:          bson.NewObjectID(),
		CreatedBy:   tempName,
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		LastReset:   *reset,
	}

	err := database.CreateCounter(c, counter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/counters/"+counter.IDHex())
}
