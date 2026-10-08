package reset

import (
	"net/http"

	"github.com/ComputerScienceHouse/counters/internal/database"
	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/ComputerScienceHouse/counters/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateReset(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var req models.CreateResetInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("failed to cast error %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	patchedResetTime := util.TranslateTime(req.ResetTime)

	if _, err := database.CreateReset(c.Request.Context(), user.Uuid, user.Username, patchedResetTime, &req); err != nil {
		logging.Logger.Warnf("failed to create error %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if _, err := database.UpdateCounterLastReset(c.Request.Context(), req.CounterID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error updating last reset")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteReset(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error getting CSHAuth reset")
		c.Status(http.StatusUnauthorized)
		return
	}

	var req models.DeleteResetInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error compiling reset")
		c.Status(http.StatusBadRequest)
		return
	}

	counterID, err := database.GetCounterFromReset(c.Request.Context(), req.RowID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error getting counter from reset")
		c.Status(http.StatusBadRequest)
		return
	}

	counterOwner, err := database.GetCounterOwner(c.Request.Context(), counterID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error getting counter owner")
		c.Status(http.StatusBadRequest)
		return
	}

	resetOwner, err := database.GetResetOwner(c.Request.Context(), req.RowID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error deleting reset")
		c.Status(http.StatusBadRequest)
		return
	}

	isCounterOwner := counterOwner == user.Uuid
	isResetOwner := resetOwner == user.Uuid

	if !users.IsEboard(user) && !users.IsActiveRTP(user) && !isCounterOwner && !isResetOwner {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := database.DeleteReset(c.Request.Context(), counterID, req.RowID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error deleting reset")
		c.Status(http.StatusInternalServerError)
		return
	}

	counterExists, err := database.UpdateCounterLastReset(c.Request.Context(), counterID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error updating last reset")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, counterExists)
}

func Routes(r *gin.RouterGroup) {
	rGroup := r.Group("/reset")

	rGroup.GET("/:id")
	rGroup.POST("", CreateReset) // create reset
	rGroup.DELETE("", DeleteReset)
}
