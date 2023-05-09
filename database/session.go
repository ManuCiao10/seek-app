package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func StoreSession(c *gin.Context, sessionID string, expiresAt time.Time) error {
	log.Println("Storing sessionID in database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client.Database("GODB").Collection("account")

	_, err := collection.UpdateOne(
		ctx, bson.M{"email": user.Email},
		bson.M{"$set": bson.M{"sessionID": sessionID, "expiresAt": expiresAt}},
	)

	if err != nil {
		log.Printf("Error updating session ID: %v", err)

		return fmt.Errorf("error updating session ID")
	}

	log.Printf("Session ID updated successfully")

	return nil
}

func CheckSession(c *gin.Context, sessionID string) error {
	log.Printf("Checking if session ID is in database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client.Database(dbName).Collection(collectionName)

	err := collection.FindOne(ctx, bson.M{"sessionID": sessionID}).Decode(&user)

	if err != nil {
		log.Printf("Session ID not found in database: %v", err)
		return fmt.Errorf("session ID not found in database")
	}

	return nil
}

func CheckSessionExpired(c *gin.Context) error {
	log.Printf("Checking if session ID is expired...")

	if time.Now().After(user.ExpiresAt) {
		log.Printf("Session ID is expired")
		return fmt.Errorf("session ID is expired")
	}

	return nil

}
