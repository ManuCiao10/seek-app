package controllers

import "github.com/gin-gonic/gin"

func GetUsers(c *gin.Context) {
	// gin.SetMode(gin.ReleaseMode)

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func HomePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
