package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"portfolio-backend/database"
	"portfolio-backend/models"
	"portfolio-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.ConnectCollection(database.Client, "user")
var validate = validator.New()

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong Password!"})
			return
		}

		passwordIsValid, msg := utils.VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Username == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		token, refreshToken, _ := utils.GenerateAllTokens(*foundUser.Name, *foundUser.Fullname, *foundUser.Username, *&foundUser.UserID)

		utils.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserID}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func SignUp() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var user models.User

			if err := ginCtx.BindJSON(&user); err != nil {
					ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
			}

			validationErr := validate.Struct(user)
			if validationErr != nil {
					ginCtx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
					return
			}

			count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
			defer cancel()
			if err != nil {
					log.Panic(err)
					ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for username"})
					return
			}

			password := utils.HashPassword(*user.Password)
			user.Password = &password

			if count > 0 {
					ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "username already exists"})
					return
			}

			user.ID = primitive.NewObjectID()
			user.UserID = user.ID.Hex()

			token, refreshToken, _ := utils.GenerateAllTokens(*user.Name, *user.Fullname, *user.Username, user.UserID)
			user.Token = &token
			user.RefreshToken = &refreshToken

			resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
			if insertErr != nil {
					msg := fmt.Sprintf("User item was not created")
					ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
					return
			}
			defer cancel()

			ginCtx.JSON(http.StatusOK, resultInsertionNumber)
	}
}

