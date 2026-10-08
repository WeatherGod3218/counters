package counters

import (
	"net/http"

	"github.com/ComputerScienceHouse/counters/internal/database"
	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

func GetCounter(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	idToUse := c.Param("id")

	counter, err := database.GetCounterFromId(c.Request.Context(), idToUse)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resets, err := database.GetResetsFromCounterId(c.Request.Context(), idToUse)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.HTML(200, "counter.tmpl", gin.H{
		"Id":          counter.CounterID,
		"CounterID":   counter.CounterOwner,
		"Title":       counter.CounterTitle,
		"Description": counter.CounterDescription,
		"Timestamp":   counter.ResetOccuredAt,
		"History":     resets,
		"Username":    user.Username,
		"FullName":    user.FullName,
		"UserID":      user.Uuid,
		"EBoard":      users.IsEboard(user),
		"RTP":         users.IsActiveRTP(user),
	})
}

func DeleteCounter(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error getting CSH auth last reset")
		c.Status(http.StatusUnauthorized)
		return
	}

	var req models.DeleteCounterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error casting request")
		c.Status(http.StatusBadRequest)
		return
	}

	counterOwner, err := database.GetCounterOwner(c.Request.Context(), req.RowID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error getting counter owner")
		c.Status(http.StatusBadRequest)
		return
	}

	sameUser := counterOwner == user.Uuid
	if !users.IsEboard(user) && !users.IsActiveRTP(user) && !sameUser {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error getting user permissions")
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := database.DeleteCounter(c.Request.Context(), req.RowID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter", "error": err}).Warning("Error deleting counter")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func CreateCounter(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var req models.CreateCounterWithResetInput

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		body, _ := c.Get(gin.BodyBytesKey)
		raw, _ := body.([]byte)

		logging.Logger.Warnf("bad request for creating counter! %s %s",
			err,
			string(raw),
		)
		c.Status(http.StatusBadRequest)
		return
	}

	counterID, err := database.CreateCounterWithReset(c.Request.Context(), user.Uuid, user.Username, &req.Counter, &req.Reset)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counters", "error": err}).Warning("Error creating new database")
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, counterID)
}

func Routes(r *gin.RouterGroup) {
	cGroup := r.Group("/counter")

	cGroup.GET("/:id")
	cGroup.POST("", CreateCounter)
	cGroup.DELETE("", DeleteCounter)
}
