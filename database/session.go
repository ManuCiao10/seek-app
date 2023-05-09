package database

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//saving session in database
//check if session is valid
//Checking if session ID is expired...

func CheckSession(c *gin.Context, sessionID string) {

	log.Printf("Checking if session ID is in database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(dbName).Collection(collectionName)

	err := collection.FindOne(ctx, bson.M{"sessionID": sessionID}).Decode(&user)

	defer cancel()

	if err != nil {
		log.Printf("Session ID not found in database: %v", err)

		c.HTML(http.StatusBadRequest, "indexnotlogged.html",
			gin.H{
				"content": "Unauthorized error: session ID not found in database",
			})

		return
	}
}

func CheckSessionExpired(c *gin.Context) {
	log.Printf("Checking if session ID is expired...")

	if time.Now().After(user.ExpiresAt) {
		log.Printf("Session ID is expired")

		c.HTML(http.StatusBadRequest, "indexnotlogged.html",
			gin.H{
				"content": "Unauthorized error: session ID is expired",
			})
		return
	}

}
