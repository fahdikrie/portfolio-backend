package utils

import (
	"context"
	"fmt"
	"log"
	"portfolio-backend/database"
	"portfolio-backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.ConnectCollection(database.Client, "user")

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Wrong password!")
		check = false
	}

	return check, msg
}

func UpdateUserByID(token string, refresh string, id string) (models.User, error) {
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: refresh})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})
	upsert := true
	filter := bson.M{"user_id": id}
	opt := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}
	opt.SetReturnDocument(options.After)
	var updatedUser models.User

	err := userCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	).Decode(&updatedUser)

	return updatedUser, err
}
