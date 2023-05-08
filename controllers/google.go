package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/seek/database"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var User_google UserGoogle
var confgoogle *oauth2.Config

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	confgoogle = &oauth2.Config{
		ClientID:     os.Getenv("ID_SECRET_GOOGLE"),
		ClientSecret: os.Getenv("CLIENT_GOOGLE"),
		RedirectURL:  "http://localhost:9000/login/google-callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func HandleGoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("HandleGoogleLogin: %v", c.Request.URL.Path)

		state, err := RandToken(32)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while generating random data."})
			return
		}

		session := sessions.Default(c)
		session.Set("state", state)
		err = session.Save()

		if err != nil {
			log.Printf("Error while saving session: %v", err)
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
			return
		}

		link := getLoginURL(state)
		c.Redirect(http.StatusTemporaryRedirect, link)

	}
}

func HandleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("HandleGoogleCallback: %v", c.Request.URL.Path)

		session := sessions.Default(c)
		retrievedState := session.Get("state")
		if retrievedState != c.Query("state") {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid session state: %s", retrievedState))
			return
		}

		tok, err := confgoogle.Exchange(context.TODO(), c.Query("code"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		client := confgoogle.Client(context.TODO(), tok)
		email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer email.Body.Close()

		data, _ := io.ReadAll(email.Body)
		err = json.Unmarshal(data, &User_google)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		log.Printf("Checking if Email is in database...")

		err = collection.FindOne(ctx, bson.M{"email": User_google.Email}).Decode(&User_google)
		defer cancel()

		if err != nil {
			log.Printf("Email %s not found in database", User_google.Email)

			_, err = collection.InsertOne(ctx, bson.M{"email": User_google.Email, "image": User_google.Picture, "name": User_google.Name, "given_name": User_google.GivenName, "family_name": User_google.FamilyName, "locale": User_google.Locale})

			if err != nil {
				log.Printf("Error inserting into database: %v", err)

				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			log.Printf("Successfully inserted into database")
			// redirect to home page

		} else {
			log.Printf("Email %s found in database", User_google.Email)
			// redirect to home page

		}

		log.Printf("Saving session %s", User_google.Email)

		session.Set("user-id", User_google.Email)
		err = session.Save()
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "[google] Successfully logged in",
		})

	}
}
