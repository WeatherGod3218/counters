package reset

import (
	"net/http"

	"github.com/ComputerScienceHouse/counters/internal/database"
	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/users"
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
		c.Status(http.StatusBadRequest)
		return
	}

	if _, err := database.CreateReset(c.Request.Context(), user.Uuid, user.Username, &req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
}

func DeleteReset(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var req models.DeleteResetInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	counterID, err := database.GetCounterFromReset(c.Request.Context(), req.RowID)
	if err := c.ShouldBindJSON(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	counterOwner, err := database.GetCounterOwner(c.Request.Context(), counterID)
	if err := c.ShouldBindJSON(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resetOwner, err := database.GetResetOwner(c.Request.Context(), req.RowID)
	if err := c.ShouldBindJSON(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	isCounterOwner := counterOwner == user.Uuid
	isResetOwner := resetOwner == user.Uuid

	if !users.IsEboard(user) && !users.IsActiveRTP(user) && !isCounterOwner && !isResetOwner {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := database.DeleteReset(c, counterID, req.RowID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error deleting reset")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func Routes(r *gin.RouterGroup) {
	rGroup := r.Group("/reset")

	rGroup.GET("/:id")
	rGroup.POST("", CreateReset) // create reset
	rGroup.DELETE("", DeleteReset)
}
