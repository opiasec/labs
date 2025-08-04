package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
}

func NewDatabase() *Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err = client.Ping(context.Background(), nil)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to MongoDB (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	fmt.Println("✅ MongoDB client initialized")

	return &Database{Client: client}
}
