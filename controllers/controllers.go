package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/seek/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// return login.html page if user not logged in
func LoginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("LoginGetHandler: %v", c.Request.URL.Path)
		sessionID, err := c.Cookie("sessionID")

		if err != nil {
			log.Printf("User is not logged in, redirect to login page")

			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}

		log.Printf("Found session ID in cookie %v", sessionID)

		err = database.CheckSession(c, sessionID)
		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{})

			return
		}

		err = database.CheckSessionExpired(c)
		if err != nil {
			if err != nil {
				c.HTML(http.StatusOK, "login.html", gin.H{})
				return
			}
		}

		log.Printf("User is logged in, redirect to index")
		log.Printf("SessionID: %v", sessionID)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})

	}
}

// Post request to login user:
// if user is valid (email and password) => store session ID in MongoDB and as a cookie
// else return login.html page with error message
func LoginPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("LoginPostHandler: %v", c.Request.URL.Path)

		if err := c.ShouldBind(&user); err != nil {
			log.Printf("Error: %v", err)

			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Invalid email or password",
					"user":    user,
				})
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

		sessionID := uuid.NewString()
		expiresAt := time.Now().Add(15 * 24 * time.Hour)

		err = database.StoreSession(c, sessionID, expiresAt)
		if err != nil {
			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Error storing session ID",
					"user":    user,
				})
			return
		}

		c.SetCookie("sessionID", sessionID, int(expiresAt.Unix()), "/", "", false, true)

		log.Printf("Session ID: %v %v", sessionID, expiresAt)
		c.Redirect(http.StatusFound, "/")

	}
}

/*
determine if the user is logged in or not by searching the sessionID in the cookie and the database MongoDB
user logged in ==> display index.html with [sign up/log in button] + [sell now button]
user not logged in ==> display index.html with [profile picture] + [sell now button] + [messages]
*/
func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("IndexGetHandler: %v", c.Request.URL.Path)

		sessionID, err := c.Cookie("sessionID")

		if len(sessionID) == 0 {
			log.Printf("User Token header is missing, user not logged in")

			c.HTML(http.StatusOK, "not_logged_in.html", gin.H{
				"content": "",
			})
			return
		}

		if err != nil {
			log.Printf("Error IndexGetHandler gettin sessionID: user is not logged in")

			c.HTML(http.StatusOK, "not_logged_in.html", gin.H{
				"content": "",
			})
			return
		}

		log.Printf("Found session ID in cookie %v", sessionID)
		err = database.CheckSession(c, sessionID)

		if err != nil {
			c.HTML(http.StatusOK, "not_logged_in.html", gin.H{})
			return
		}

		err = database.CheckSessionExpired(c)

		if err != nil {
			c.HTML(http.StatusOK, "not_logged_in.html", gin.H{})
			return
		}

		log.Printf("User logged in, session ID: %v", sessionID)
		c.Redirect(http.StatusFound, "/")
	}
}

// Sign up user with the salted password and save user in MongoDB database, then redirect to login page
// If user already exists, return login.html page
func SignupPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("SignupPostHandler: %v", c.Request.URL.Path)

		if err := c.ShouldBind(&snUser); err != nil {
			log.Printf("Error struct: %v", err)

			c.HTML(http.StatusBadRequest, "signup.html",
				gin.H{
					"content": "Invalid email or password",
					"user":    user,
				})
			return
		}
		log.Printf("User: %v", snUser)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(snUser.Password), 8)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		err = collection.FindOne(ctx, bson.M{"email": snUser.Email}).Decode(&snUser)

		defer cancel()

		if err == nil {
			log.Printf("User already exists: %v", err)

			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "User already exists",
					"user":    user,
				})
			return
		}

		_, err = collection.InsertOne(ctx, bson.M{"email": snUser.Email, "password": string(hashedPassword), "fullname": snUser.Fullname, "username": snUser.Username})
		if err != nil {
			log.Printf("Error: %v", err)
		}
		log.Printf("User created: %v", snUser.Email)

		c.Redirect(http.StatusFound, "/login")

	}
}

// Display the signup.html page
func SignupGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("SignupGetHandler: %v", c.Request.URL.Path)

		c.HTML(http.StatusOK, "signup.html", gin.H{
			"content": "",
		})
	}
}
