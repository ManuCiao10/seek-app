package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seek/docs/database"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/seek/docs/controllers"
)

// AuthRequired middleware to check if user is logged in
// if user is logged in => redirect to index page
// if user is not logged in => redirect to login page
func AuthRequired(c *gin.Context) {
	log.Printf("AuthRequired: %v", c.Request.URL.Path)

	sessionID, err := c.Cookie("sessionID")

	if len(sessionID) == 0 {
		log.Printf("User Token header is missing, redirect to index page with [sign up/log in button] + [sell now button]")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "",
		})
		return
	}

	if err != nil {
		log.Printf("User is not logged in, redirect to index page with [sign up/log in button] + [sell now button]")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "",
		})
		return
	}

	log.Printf("Found session ID in cookie %v", sessionID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.Client.Database("GODB").Collection("account")

	log.Printf("Checking if session ID is in database...")

	err = collection.FindOne(ctx, bson.M{"sessionID": sessionID}).Decode(&controllers.AuthUser)

	defer cancel()

	if err != nil {
		log.Printf("Session ID not found in database: %v", err)

		c.HTML(http.StatusBadRequest, "index.html",
			gin.H{
				"content": "Unauthorized error: session ID not found in database",
			})
		return
	}

	log.Printf("Checking if session ID is expired...")

	if time.Now().After(controllers.AuthUser.ExpiresAt) {
		log.Printf("Session ID is expired")

		c.HTML(http.StatusBadRequest, "index.html",
			gin.H{
				"content": "Unauthorized error: session ID is expired",
			})
		return
	}

	log.Printf("User is logged in, redirect to index")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": "",
	})

}
