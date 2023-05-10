package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

//checking if user already exist

func CheckEmail(c *gin.Context, email string) error {
	log.Printf("Checking if Email is in database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database("GODB").Collection("account")

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&UserDB)

	defer cancel()

	if err != nil {
		return fmt.Errorf("error Email")
	}

	return nil
}

func CheckPassword(c *gin.Context, userPassword string) error {
	userPass := []byte(userPassword)
	dbPass := []byte(UserDB.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	if passErr != nil {
		log.Printf("Error Password: %s", passErr)

		return fmt.Errorf("error Password")
	}

	return nil
}
