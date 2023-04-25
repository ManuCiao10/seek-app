package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/seek/docs/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	user     UserPostLogin  // user from login form
	dbUser   UserPostLogin  // user from database
	snUser   UserPostSignup // user from signup form
	AuthUser UserPostSignup // user from database AuthRequired middleware
)

// return login.html page if user not logged in
func LoginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("LoginGetHandler: %v", c.Request.URL.Path)
		sessionID, err := c.Cookie("sessionID")

		if err != nil {
			log.Printf("User is not logged in, redirect to login page")

			c.HTML(http.StatusOK, "login.html", gin.H{
				"content": "",
			})
			return
		}

		log.Printf("Found session ID in cookie %v", sessionID)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		log.Printf("Checking if session ID is in database...")

		err = collection.FindOne(ctx, bson.M{"sessionID": sessionID}).Decode(&snUser)

		defer cancel()

		if err != nil {
			log.Printf("Session ID not found in database: %v", err)

			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Unauthorized error: session ID not found in database",
				})
			return
		}

		log.Printf("Checking if session ID is expired...")

		if time.Now().After(snUser.ExpiresAt) {
			log.Printf("Session ID is expired")

			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Unauthorized error: session ID is expired",
				})
			return
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

		log.Println("User is valid:")
		log.Println("Storing sessionID in database...")

		sessionID := uuid.NewString()
		expiresAt := time.Now().Add(15 * 24 * time.Hour)

		_, err = collection.UpdateOne(
			ctx, bson.M{"email": user.Email},
			bson.M{"$set": bson.M{"sessionID": sessionID, "expiresAt": expiresAt}},
		)

		if err != nil {
			log.Printf("Error updating session ID: %v", err)

			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Error updating session ID",
					"user":    user,
				})
			return
		}

		c.SetCookie("sessionID", sessionID, int(expiresAt.Unix()), "/", "", false, true)

		log.Printf("Session ID: %v %v", sessionID, expiresAt)
		c.Redirect(http.StatusFound, "/")

	}
}

// display the index.html page
// determine if the user is logged in or not by searching the sessionID in the cookie and the database MongoDB
// user logged in ==> display index.html with [sign up/log in button] + [sell now button]
// user not logged in ==> display index.html with [profile picture] + [sell now button] + [messages]

func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("IndexGetHandler: %v", c.Request.URL.Path)

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

		err = collection.FindOne(ctx, bson.M{"sessionID": sessionID}).Decode(&snUser)

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

		if time.Now().After(snUser.ExpiresAt) {
			log.Printf("Session ID is expired")

			c.HTML(http.StatusBadRequest, "index.html",
				gin.H{
					"content": "Unauthorized error: session ID is expired",
				})
			return
		}

		log.Printf("User is logged in, redirect to index")
		log.Printf("Session ID: %v", sessionID)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})

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
