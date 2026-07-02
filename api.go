package main

import (
	"net/http"
	"sort"

	"github.com/WeatherGod3218/counters/database"
	"github.com/WeatherGod3218/counters/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	cshAuth "github.com/computersciencehouse/csh-auth/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetHomePage(c *gin.Context) {
	userAny, exists := c.Get("cshauth")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}
	user, ok := userAny.(*cshAuth.Claims)
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

func GetCreatePage(c *gin.Context) {
	userAny, exists := c.Get("cshauth")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}
	user, ok := userAny.(*cshAuth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}

	c.HTML(http.StatusOK, "create.tmpl", gin.H{
		"Username": user.Username,
		"FullName": user.FullName,
		"EBoard":   IsEboard(user),
		"RTP":      IsActiveRTP(user),
	})
}

func GetResetPage(c *gin.Context) {
	userAny, exists := c.Get("cshauth")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}
	user, ok := userAny.(*cshAuth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}

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
	userAny, exists := c.Get("cshauth")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}
	user, ok := userAny.(*cshAuth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}
	idToUse := c.Param("id")

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
	userAny, exists := c.Get("cshauth")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}
	user, ok := userAny.(*cshAuth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}

	resetTime := TranslateTime(c.PostForm("reset-time"))

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
	userAny, exists := c.Get("cshauth")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}
	user, ok := userAny.(*cshAuth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load CSH context!"})
		return
	}

	resetTime := TranslateTime(c.PostForm("reset-time"))

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
