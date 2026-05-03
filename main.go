package main

import (
	"net/http"

	//cshAuth "github.com/computersciencehouse/csh-auth"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")
	// csh := cshAuth.CSHAuth{}

	// csh.Init(
	// 	os.Getenv("AUTH_OIDC_ID"),
	// 	os.Getenv("AUTH_OIDC_SECRET"),
	// 	os.Getenv("AUTH_JWC_SECRET"),
	// 	os.Getenv("VOTE_STATE"),
	// 	os.Getenv("SERVER_HOST"),
	// 	os.Getenv("SERVER_HOST")+"/auth/callback",
	// 	os.Getenv("SERVER_HOST")+"/auth/login",
	// 	[]string{"profile", "email", "groups"},
	// )

	// router.GET("/auth/login", csh.AuthRequest)
	// router.GET("/auth/callback", csh.AuthCallback)
	// router.GET("/auth/logout", csh.AuthLogout)

	//meow

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "My Site",
		})
	})

	router.Run(":8080")
}
