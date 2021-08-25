package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// server for testing the webhook demo on the eventbus
func main() {
	router := gin.Default()

	router.POST("/stream/:name", func(c *gin.Context) {
		name := c.Param("uuid")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.Run(":8080")
}