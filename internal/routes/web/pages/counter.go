package pages

import (
	"net/http"

	"github.com/ComputerScienceHouse/counters/internal/database"
	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetCounterPage(c *gin.Context) {
	user, err := users.GetCSHAuth(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	idToUse := c.Param("id")

	counter, err := database.GetCounterFromId(c.Request.Context(), idToUse)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "counter page", "error": err}).Warning("Error geting counter")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse())
		return
	}

	if counter == nil {
		c.Redirect(http.StatusNotFound, "/")
		return
	}

	history, err := database.GetResetsFromCounterId(c.Request.Context(), idToUse)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "reset", "error": err}).Warning("Error getting history for counter")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse())
		return
	}

	logging.Logger.Info(counter.CounterDescription)

	c.HTML(http.StatusOK, "counter.tmpl", gin.H{
		"Id":          counter.CounterID,
		"CounterID":   counter.CounterOwner,
		"Title":       counter.CounterTitle,
		"Description": counter.CounterDescription,
		"Timestamp":   counter.ResetOccuredAt,

		"History":  history,
		"Username": user.Username,
		"FullName": user.FullName,
		"UserID":   user.Uuid,

		"EBoard": users.IsEboard(user),
		"RTP":    users.IsActiveRTP(user),
	})
}
