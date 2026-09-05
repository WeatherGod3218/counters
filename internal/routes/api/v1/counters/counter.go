package counters

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	cGroup := r.Group("/counter")

	cGroup.GET("/:id")
	cGroup.POST("")
	cGroup.PATCH("")
	cGroup.DELETE("")
}
