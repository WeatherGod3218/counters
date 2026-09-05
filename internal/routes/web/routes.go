package pages

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	r.GET("/")
	r.GET("/create")
	r.GET("/reset/:id")
}
