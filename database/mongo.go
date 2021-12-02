package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = MongoDBInstance()

func MongoDBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	MongoDBServer := os.Getenv("MONGODB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDBServer))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func ConnectCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	MongoDBCluster := os.Getenv("MONGODB_CLUSTER")
	var collection *mongo.Collection = client.Database(MongoDBCluster).Collection(collectionName)
	return collection
}
