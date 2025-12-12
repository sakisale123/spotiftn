package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var MongoDB *mongo.Database

func ConnectMongo() {
	// context sa timeout-om
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal("Mongo connection error: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Mongo ping failed: ", err)
	}

	MongoDB = client.Database("authdb")
	UserCollection = MongoDB.Collection("users")

	log.Println("Connected to MongoDB (users service)")
}
