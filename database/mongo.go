package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitMongoDB() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable")
	}

	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")

	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	log.Println("Ping to MongoDB!")
}

func DeleteUser(id string) error {
	collection := Client.Database("GODB").Collection("account")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		return err
	}

	return nil
}
