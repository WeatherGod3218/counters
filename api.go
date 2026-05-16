package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/WeatherGod3218/counters/database"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/v2/bson"
)

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
	user := GetUserData(c)

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
		"Username": user.Username,
		"FullName": user.FullName,
		"EBoard":   IsEboard(user),
		"RTP":      IsActiveRTP(user),
	})
}

func GetCreatePage(c *gin.Context) {
	user := GetUserData(c)

	c.HTML(http.StatusOK, "create.tmpl", gin.H{
		"Username": user.Username,
		"FullName": user.FullName,
		"EBoard":   IsEboard(user),
		"RTP":      IsActiveRTP(user),
	})
}

func GetResetPage(c *gin.Context) {
	user := GetUserData(c)

	idToUse := c.Param("id")
	c.HTML(http.StatusOK, "reset.tmpl", gin.H{
		"Id":       idToUse,
		"Username": user.Username,
		"FullName": user.FullName,
		"EBoard":   IsEboard(user),
		"RTP":      IsActiveRTP(user),
	})
}

func LoadCounter(c *gin.Context) {
	idToUse := c.Param("id")
	user := GetUserData(c)

	counter, err := database.GetCounterFromId(c, idToUse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	history, err := database.GetHistoryFromId(c, idToUse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "counter.tmpl", gin.H{
		"Id":          counter.Id.Hex(),
		"CounterID":   counter.UserID,
		"Title":       counter.Title,
		"Description": counter.Description,
		"Timestamp":   counter.LastReset.Timestamp,
		"History":     history.History,
		"Username":    user.Username,
		"FullName":    user.FullName,
		"UserID":      user.Uuid,
		"EBoard":      IsEboard(user),
		"RTP":         IsActiveRTP(user),
	})
}

func ResetCounter(c *gin.Context) {
	resetTime := translateTime(c.PostForm("reset-time"))
	user := GetUserData(c)

	reset := &database.Reset{
		Id:          bson.NewObjectID(),
		UserID:      user.Uuid,
		Reporter:    user.Username,
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
	user := GetUserData(c)

	reset := &database.Reset{
		Id:          bson.NewObjectID(),
		UserID:      user.Uuid,
		Reporter:    user.Username,
		Description: c.PostForm("reset-description"),
		Timestamp:   resetTime,
	}

	counter := &database.Counter{
		Id:          bson.NewObjectID(),
		UserID:      user.Uuid,
		CreatedBy:   user.Username,
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

func DeleteReset(c *gin.Context) {
	counterIdString := c.PostForm("counterId")
	resetId, _ := bson.ObjectIDFromHex(c.PostForm("resetId"))
	counterId, _ := bson.ObjectIDFromHex(counterIdString)

	exists, err := database.DeleteReset(c, counterId, resetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.Redirect(http.StatusFound, "/counters/"+counterIdString)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func DeleteCounter(c *gin.Context) {
	counterId, _ := bson.ObjectIDFromHex(c.PostForm("counterId"))

	err := database.DeleteCounter(c, counterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
