package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seek/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Sign up user with the salted password and save user in MongoDB database, then redirect to login page
// If user already exists, return login.html page
func SignupPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("SignupPostHandler: %v", c.Request.URL.Path)

		fullname := c.PostForm("fullname")
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		err = collection.FindOne(ctx, bson.M{"email": email}).Decode(&database.UserDB)

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

		log.Printf("adding user to database: %v", email)
		_, err = collection.InsertOne(ctx, bson.M{"email": email, "password": string(hashedPassword), "fullname": fullname, "username": username})
		if err != nil {
			log.Printf("Error adding user to database: %v", err)

			c.HTML(http.StatusBadRequest, "signup.html",
				gin.H{
					"content": "Error adding user to database",
					"user":    user,
				})
			return

		}
		log.Printf("User created: %v", email)

		c.Redirect(http.StatusFound, "/login")

	}
}

func SignupGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("SignupGetHandler: %v", c.Request.URL.Path)

		c.HTML(http.StatusOK, "signup.html", gin.H{})
	}
}
