package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	init_mongo()
}

func init_mongo() {

	fmt.Println("initializing MongoDB...")
	// Set client options
	clientOptions := options.Client().ApplyURI(MONGODB_URI + ":" + MONGODB_PORT)

	// Connect to MongoDB
	client, _ = mongo.Connect(context.Background(), clientOptions)
	// Check the connection
	err := client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Unable to connect MongoDB ", err)
	}

	fmt.Println("Connected to MongoDB.")
}
