package main

import (
	"net/http"
	"os"

	"github.com/ComputerScienceHouse/counters/internal/database"
	v1 "github.com/ComputerScienceHouse/counters/internal/routes/api/v1"
	"github.com/ComputerScienceHouse/counters/internal/routes/web"
	"github.com/gin-gonic/gin"

	"github.com/ComputerScienceHouse/counters/internal/logging"
	"github.com/sirupsen/logrus"

	cshAuth "github.com/computersciencehouse/csh-auth/v2"
)

var DEV_FORCE_IS_EBOARD bool = os.Getenv("DEV_FORCE_IS_EBOARD") == "true"

func main() {

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

	if err := database.InitDatabase(); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error initializing database")
	}
	router := gin.Default()

	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")

	router.GET("/auth/login", auth.HandleLogin)
	router.GET("/auth/callback", auth.HandleCallback)
	router.GET("/auth/logout", auth.HandleLogout)

	router.Use(auth.CookieMiddleware())

	web.Routes(router.Group(""))
	v1.Routes(router.Group(""))

	router.Run(":8080")
}
