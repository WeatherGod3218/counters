package main

import (
	"net/http"
	"os"

	"github.com/WeatherGod3218/counters/database"
	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/counters/logging"
	"github.com/sirupsen/logrus"

	cshAuth "github.com/computersciencehouse/csh-auth/v2"
)

var DEV_FORCE_IS_EBOARD bool = os.Getenv("DEV_FORCE_IS_EBOARD") == "true"

func main() {
	database.Client = database.Connect()

	hostUrl := os.Getenv("SERVER_HOST")
	auth, err := cshAuth.Init(
		os.Getenv("AUTH_OIDC_ID"),
		os.Getenv("AUTH_OIDC_SECRET"),
		hostUrl,
		hostUrl+"/auth/login",
		hostUrl+"/auth/callback",
		[]string{"profile", "email", "groups"},
	)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error initializing csh-auth")
	}

	router := gin.Default()

	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")

	router.GET("/auth/login", auth.HandleLogin)
	router.GET("/auth/callback", auth.HandleCallback)
	router.GET("/auth/logout", auth.HandleLogout)

	router.Use(auth.CookieMiddleware())

	router.GET("/", GetHomePage)
	router.GET("/counters/:id", LoadCounter)

	router.GET("/create", GetCreatePage)
	router.POST("/create", CreateCounter)

	router.GET("/reset/:id", GetResetPage)
	router.POST("/reset/:id", ResetCounter)

	router.POST("/delete/counter", DeleteCounter)
	router.POST("/delete/reset", DeleteReset)

	router.Run(":8080")
}
