package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/seek/docs/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	user      UserPostLogin // user from login form
	dbUser    UserPostLogin // user from database
	validUser = true        // check if user is valid
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

		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)

		defer cancel()

		if err != nil {
			log.Printf("Error Email: %v", err)
			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Invalid email ",
					"user":    user,
				})
			return
		}

		userPass := []byte(user.Password)
		dbPass := []byte(dbUser.Password)

		passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
		if passErr != nil {
			log.Printf("Error Password: %v", err)

			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Invalid password",
					"user":    user,
				})
			return
		}

		log.Println("User is valid")

		sessionID := uuid.New().String()

		// Store session ID in Redis store
		session := sessions.Default(c)
		session.Set("sessionID", sessionID)
		session.Save()

		// Set session ID as a cookie
		c.SetCookie("sessionID", sessionID, 3600, "/", "", false, true)

		c.Redirect(http.StatusFound, "/")

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

// TODO: sign up user with the salted password
