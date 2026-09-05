package reset

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	rGroup := r.Group("/resets")

	rGroup.GET("/:id")

	rGroup.POST("")   // create reset
	rGroup.DELETE("") // delete reset
}
