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

var oidcClient = OIDCClient{}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("claims")

		if !exists {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			return
		}

		claims, ok := value.(*cshAuth.Claims)

		if !ok || claims == nil {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Set("authClaims", claims)

		c.Next()
	}
}

func main() {
	database.Client = database.Connect()

	auth, err := cshAuth.Init(
		os.Getenv("AUTH_OIDC_ID"),
		os.Getenv("AUTH_OIDC_SECRET"),
		os.Getenv("SERVER_HOST"),
		os.Getenv("SERVER_HOST")+"/auth/login",
		os.Getenv("SERVER_HOST")+"/auth/callback",
		[]string{"profile", "email", "groups"},
	)
	oidcClient.setupOidcClient(os.Getenv("AUTH_OIDC_ID"), os.Getenv("AUTH_OIDC_SECRET"))

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error initializing csh-auth")
	}

	router := gin.Default()

	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")

	router.GET("/auth/login", auth.HandleLogin)       // This endpoint should match the path for loginURL
	router.GET("/auth/callback", auth.HandleCallback) // This endpoint should match the path for callbackURL
	router.GET("/auth/logout", auth.HandleLogout)

	router.Use(auth.CookieMiddleware())
	router.Use(RequireAuth())

	router.GET("/", GetHomePage)
	router.GET("/counters/:id", LoadCounter)

	router.GET("/create", GetCreatePage)
	router.POST("/create", CreateCounter)

	router.GET("/reset/:id", GetResetPage)
	router.POST("/reset/:id", ResetCounter)

	router.DELETE("/delete/counter", DeleteCounter)
	router.DELETE("/delete/reset", DeleteReset)

	router.Run(":8080")
}
