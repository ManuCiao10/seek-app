package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/seek/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	confdiscord = &oauth2.Config{
		RedirectURL: "http://localhost:9000/login/discord-callback",

		ClientID:     os.Getenv("ID_SECRET_DISCORD"),
		ClientSecret: os.Getenv("CLIENT_DISCORD"),
		Scopes: []string{
			"guilds",
			"identify",
			"email",
		},

		Endpoint: Endpoint,
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

		link := getDiscordLoginURL(state)
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

		token, err := confdiscord.Exchange(context.TODO(), c.Query("code"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		client := confdiscord.Client(context.TODO(), token)
		resp, err := client.Get("https://discord.com/api/users/@me")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)

			return
		}

		// fmt.Printf("Discord Response: %s\n", body) // DEBUG

		err = json.Unmarshal(body, &User_discord)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		log.Printf("Checking if email [%s] is in database...", User_discord.Email)

		avatar_image := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", User_discord.ID, User_discord.Avatar)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		collection := database.Client.Database("GODB").Collection("account")

		err = collection.FindOne(ctx, bson.M{"email": User_discord.Email}).Decode(&User_discord)
		defer cancel()

		if err != nil {
			log.Printf("Email %s not found in database", User_discord.Email)

			_, err = collection.InsertOne(ctx, bson.M{"email": User_discord.Email, "image": avatar_image, "username_discord": User_discord.Username, "given_name": "to_fix", "family_name": "to_fix", "locale": User_discord.Locale})

			if err != nil {
				log.Printf("Error inserting into database: %v", err)

				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			log.Printf("Successfully inserted into database")

		} else {
			log.Printf("Email %s found in database", User_discord.Email)

		}

		log.Printf("Saving session %s", User_discord.Email)

		session.Set("user-id", User_discord.Email)
		err = session.Save()
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
			return
		}

		log.Println("Storing sessionID in database...")

		sessionID := uuid.NewString()
		expiresAt := time.Now().Add(15 * 24 * time.Hour)

		err = database.StoreSession(c, sessionID, expiresAt, User_discord.Email)
		if err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{})
			return
		}

		c.SetCookie("sessionID", sessionID, int(expiresAt.Unix()), "/", "", false, true)

		log.Printf("Session ID: %v %v", sessionID, expiresAt)
		c.Redirect(http.StatusFound, "/")

	}
}
