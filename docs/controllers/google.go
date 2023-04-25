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
	"github.com/seek/docs/database"
	"go.mongodb.org/mongo-driver/bson"

	// "github.com/golang/glog"

	// "github.com/golang/glog"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	// "google.golang.org/api/option"
)

const (
	stateKey  = "state"
	sessionID = "ginoauth_google_session"
)

var User UserGoogle
var conf *oauth2.Config

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf = &oauth2.Config{
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
			log.Printf("Error while generating random data: %v", err)
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

		tok, err := conf.Exchange(context.TODO(), c.Query("code"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		client := conf.Client(context.TODO(), tok)
		email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer email.Body.Close()

		data, _ := io.ReadAll(email.Body)
		err = json.Unmarshal(data, &User)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		log.Printf("Checking if Email is in database...")

		err = collection.FindOne(ctx, bson.M{"email": User.Email}).Decode(&User)
		defer cancel()

		if err != nil {
			log.Printf("Email %s not found in database", User.Email)

			_, err = collection.InsertOne(ctx, bson.M{"email": User.Email, "image": User.Picture, "name": User.Name, "given_name": User.GivenName, "family_name": User.FamilyName, "locale": User.Locale})

			if err != nil {

				log.Printf("Error inserting into database: %v", err)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			log.Printf("Successfully inserted into database")
			// redirect to home page

		} else {
			log.Printf("Email %s found in database", User.Email)
			// redirect to home page

		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully logged in",
		})

	}
}
