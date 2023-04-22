package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/seek/docs/database"
	"go.mongodb.org/mongo-driver/bson"
)

// return login.html page if user not logged in
func LoginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("sessionID")

		if err != nil {
			log.Printf("User is not logged in, redirect to login page")

			c.HTML(http.StatusOK, "login.html", gin.H{
				"content": "",
			})
			return
		}

		session := sessions.Default(c)
		sessionIDFromStore := session.Get("sessionID")
		if sessionIDFromStore == nil || sessionIDFromStore.(string) != sessionID {
			log.Printf("Session ID in cookie does not match session ID in Redis store")

			c.HTML(http.StatusOK, "login.html", gin.H{
				"content": "",
			})
			return
		}

		log.Printf("User is logged in, redirect to index")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})

	}
}

// Post request to login user and save session
func LoginPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UserPostLogin
		var dbUser UserPostLogin
		var validUser bool

		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		fmt.Print(user.Email)
		err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)

		if err != nil {
			log.Printf("Error MongoDB: %v", err)
			validUser = false
		}

		defer cancel()

		log.Printf("User from DB: %v", dbUser)

		if validUser {
			log.Println("User is valid")

			sessionID := uuid.New().String()

			// Store session ID in Redis store
			session := sessions.Default(c)
			session.Set("sessionID", sessionID)
			session.Save()

			// Set session ID as a cookie
			c.SetCookie("sessionID", sessionID, 3600, "/", "", false, true)

			c.Redirect(http.StatusFound, "/")
		} else {
			log.Println("User is invalid")
			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Invalid email or password",
					"user":    user,
				})

		}

	}
}

// return index.html page (if user not logged in dipaly index.html page [with the login button])
// if user logged in display index.html page
func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve session ID from cookie
		sessionID, err := c.Cookie("sessionID")
		if err != nil {
			// User is not logged in, redirect to login page
			c.HTML(http.StatusOK, "indexNoLogin.html", gin.H{
				"content": "This is an index page...",
			})
			return
		}

		// Load session data from Redis store
		session := sessions.Default(c)
		sessionIDFromStore := session.Get("sessionID")
		if sessionIDFromStore == nil || sessionIDFromStore.(string) != sessionID {
			// Session ID in cookie does not match session ID in Redis store
			c.HTML(http.StatusOK, "indexNoLogin.html", gin.H{
				"content": "This is an index page...",
			})
			return
		}

		// User is logged in, render main page (without login button)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})
	}
}

// Salt and hash the password using the bcrypt algorithm
