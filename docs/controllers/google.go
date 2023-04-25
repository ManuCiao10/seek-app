package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seek/docs/database"
	"github.com/seek/docs/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserGoogle struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func HandleGoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("HandleGoogleLogin: %v", c.Request.URL.Path)

		conf := &oauth2.Config{
			ClientID:     utils.GetEnv("ID_SECRET_GOOGLE", ""),
			ClientSecret: utils.GetEnv("CLIENT_GOOGLE", ""),
			RedirectURL:  "http://localhost:9000/login/google-callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}

		url := conf.AuthCodeURL("state")
		c.Redirect(http.StatusTemporaryRedirect, url)

	}
}

func HandleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := &oauth2.Config{
			ClientID:     utils.GetEnv("ID_SECRET_GOOGLE", ""),
			ClientSecret: utils.GetEnv("CLIENT_GOOGLE", ""),
			RedirectURL:  "http://localhost:9000/login/google-callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}

		// session := sessions.Default(c)
		// retrievedState := session.Get("state")
		// if retrievedState != c.Query("state") {
		// 	c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		// 	return
		// }

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

		//save data to struct and save to database
		var UserGoogle UserGoogle

		data, _ := io.ReadAll(email.Body)
		err = json.Unmarshal(data, &UserGoogle)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		log.Printf("Checking if Email is in database...")

		err = collection.FindOne(ctx, bson.M{"email": UserGoogle.Email}).Decode(&UserGoogle)
		defer cancel()

		if err != nil {
			log.Printf("Email %s not found in database", UserGoogle.Email)
			// redi

			_, err = collection.InsertOne(ctx, bson.M{"email": UserGoogle.Email, "image": UserGoogle.Picture, "name": UserGoogle.Name, "given_name": UserGoogle.GivenName, "family_name": UserGoogle.FamilyName, "locale": UserGoogle.Locale})

			if err != nil {

				log.Printf("Error inserting into database: %v", err)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			log.Printf("Successfully inserted into database")
			// redirect to home page

		} else {
			log.Printf("Email %s found in database", UserGoogle.Email)
			// redirect to home page

		}

		c.Status(http.StatusOK)

	}
}
