package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo() *mongo.Client {
	envVariables := EnvVariable

	fmt.Println("connecting to mongo...")
	client, err := mongo.NewClient(options.Client().ApplyURI(envVariables.dbUri))
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cancel() // Need to understand context!!

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to mongo")
	return client
}

var DB *mongo.Client = ConnectToMongo()

func GetCollection(db *mongo.Client, collectionName string) *mongo.Collection {
	collection := db.Database("golangDb").Collection(collectionName)
	return collection
}
