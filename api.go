package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomePage(c *gin.Context) {
	// This is intentionally left unprotected
	// A user may be unable to vote but should still be able to see a list of polls

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Counters": "Hello",
	})
}

func GetCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "create.tmpl", gin.H{
		"No": "Yes",
	})
}
