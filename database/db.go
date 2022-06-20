package database

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func InitCollection(name string) *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("error creating database client")
	}

	db := client.Database("notification-service")

	return db.Collection(name)
}
