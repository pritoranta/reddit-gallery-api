package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	router.GET("/subreddits", QuerySubredditsEndpoint)
	router.Run()
}
