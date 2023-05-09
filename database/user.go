package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

//checking if user already exist
//Checking if Email is in database...

func DeleteUser(id string) error {
	collection := Client.Database("GODB").Collection("account")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		return err
	}

	return nil
}
