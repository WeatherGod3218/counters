package web

import (
	"github.com/ComputerScienceHouse/counters/internal/routes/web/pages"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.GET("/", pages.GetHomePage)
	r.GET("/create", pages.GetCreatePage)
	r.GET("/reset/:id", pages.GetResetPage)
	r.GET("/counter/:id", pages.GetCounterPage)
}
