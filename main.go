package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	// Set up the client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Print message when connection is successful
	fmt.Println("Connected to MongoDB successfully")

	// Get a handle to the database
	db := client.Database("sokuja")

	// Get a handle to the collection
	collection := db.Collection("Token")

	// Update the existing document every 10 minutes
	objectID, err := primitive.ObjectIDFromHex("6676d1e5c08d4d036bcff8c6")
	if err != nil {
		log.Fatalf("failed to convert ObjectID: %v", err)
	}

	for {
		// Update the existing document
		filter := bson.M{"_id": objectID}
		update := bson.M{"$set": bson.M{"token": generateRandomString(87)}}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatalf("failed to update document: %v", err)
		}

		fmt.Println("Document updated successfully")

		// Wait for 10 minutes
		time.Sleep(10 * time.Minute)
	}
}
