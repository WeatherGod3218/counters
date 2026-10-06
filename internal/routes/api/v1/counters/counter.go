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

func DeleteCounter(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var req models.DeleteCounterInput
	if err := c.ShouldBindJSON(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	counterOwner, err := database.GetCounterOwner(c.Request.Context(), req.RowID)
	if err := c.ShouldBindJSON(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	sameUser := counterOwner == user.Uuid
	if !users.IsEboard(user) && !users.IsActiveRTP(user) && !sameUser {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := database.DeleteCounter(c.Request.Context(), req.RowID); err != nil {
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

	if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
		body, _ := c.Get(gin.BodyBytesKey)
		raw, _ := body.([]byte)

		logging.Logger.Warnf("bad request for creating counter! %s %s",
			err,
			string(raw),
		)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := database.CreateCounterWithReset(c.Request.Context(), user.Uuid, user.Username, &req.Counter, &req.Reset); err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counters", "error": err}).Warning("Error creating new database")
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusNoContent)
}

func Routes(r *gin.RouterGroup) {
	cGroup := r.Group("/counter")

	cGroup.GET("/:id")
	cGroup.POST("", CreateCounter)
	cGroup.DELETE("", DeleteCounter)
}
