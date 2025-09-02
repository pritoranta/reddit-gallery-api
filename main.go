package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://pritoranta.net", "http://localhost:5173"},
	}))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	router.GET("/subreddits", QuerySubredditsEndpoint)
	router.GET("/media", QueryMediaEndpoint)
	router.Run(":9361")
}
