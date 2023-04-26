package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var confgoogle *oauth2.Config

// Endpoint is Discord's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://discord.com/api/oauth2/authorize",
	TokenURL:  "https://discord.com/api/oauth2/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	confgoogle = &oauth2.Config{
		RedirectURL: "http://localhost:9000/login/discord-callback",
		// This next 2 lines must be edited before running this.
		ClientID:     os.Getenv("ID_SECRET_DISCORD"),
		ClientSecret: os.Getenv("CLIENT_DISCORD"),
		Scopes:       []string{ScopeIdentify},
		Endpoint:     Endpoint,
	}

}

func HandleDiscordLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("HandleDiscordLogin: %v", c.Request.URL.Path)

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

func HandleDiscordCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("HandleDiscordCallback: %v", c.Request.URL.Path)

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
		token, err := client.Get("https://discord.com/api/users/@me")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer token.Body.Close()

		log.Printf("Token: %v", token)

	}
}
