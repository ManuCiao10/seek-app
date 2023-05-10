package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seek/database"
)

/*
Session Destruction: When a user logs out or their session expires,
you will need to destroy the session.
This typically involves removing the session data from the session store.
*/
func LogoutGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("LogoutGetHandler: %v", c.Request.URL.Path)

		c.SetCookie("sessionID", "", -1, "/", "", false, true)

		err := database.DeleteSession(c)
		if err != nil {
			log.Printf("LogoutGetHandler: %v", err)
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{})
			return
		}

		log.Printf("LogoutGetHandler: [%v]", "User logged out")
		c.Redirect(http.StatusFound, "/")
	}
}
