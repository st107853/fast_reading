package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const dbName = "library"
const collName = "books"

type Logger struct {
	*mongo.Client // The database access interface
}

var db Logger

func Connect() error {
	// Capture connection propeties.
	//	host := os.Getenv("HOST")

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	db = Logger{client}

	fmt.Println("Connected to MongoDB!")
	return nil
}
