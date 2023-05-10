package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seek/database"
)

// if user is not logged in: redirect to login page
func AuthRequired(c *gin.Context) {
	log.Printf("AuthRequired: %v", c.Request.URL.Path)

	sessionID, err := c.Cookie("sessionID")

	if len(sessionID) == 0 {
		log.Printf("User Token header is missing, redirect to login page")

		c.Redirect(http.StatusFound, "/login")
		return
	}

	if err != nil {
		log.Printf("Error getting session ID from cookie: %v", err)

		c.Redirect(http.StatusFound, "/login")
		return
	}

	log.Printf("Found session ID in cookie %v", sessionID)

	err = database.CheckSession(c, sessionID)
	if err != nil {
		log.Printf("Session ID not found in database: %v", err)

		c.Redirect(http.StatusFound, "/login")
		return
	}

	err = database.CheckSessionExpired(c)
	if err != nil {
		log.Printf("Session ID is expired: %v", err)

		c.Redirect(http.StatusFound, "/login")
		return
	}

	log.Printf("User is logged in")
	c.Next()
}
