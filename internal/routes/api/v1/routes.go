package v1

import (
	"github.com/WeatherGod3218/counters/internal/routes/api/v1/counters"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	v1Group := r.Group("/v1")

	counters.Routes(v1Group)
}
