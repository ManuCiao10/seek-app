package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
)

/*
Session Destruction: When a user logs out or their session expires, 
you will need to destroy the session. 
This typically involves removing the session data from the session store.
*/
func LogoutGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("LogoutGetHandler: %v", c.Request.URL.Path)

		// delete cookie

	}
}
