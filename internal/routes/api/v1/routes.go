package v1

import (
	"github.com/ComputerScienceHouse/counters/internal/routes/api"
	"github.com/ComputerScienceHouse/counters/internal/routes/api/v1/counters"
	"github.com/ComputerScienceHouse/counters/internal/routes/api/v1/reset"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	v1Group := r.Group("/api/v1", api.CookieToHeaderAuth())

	counters.Routes(v1Group)
	reset.Routes(v1Group)
}
