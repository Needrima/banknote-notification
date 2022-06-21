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

// InitCollection intitalizes the database and returned the collection with name "name"
func InitCollection(name string) *mongo.Collection {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading environmental variables")
	}

	context, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_CONN_STR"))

	client, err := mongo.Connect(context, clientOptions)
	if err != nil {
		log.Fatal("error creating database client:", err)
	} else {
		log.Println("connection successful")
	}

	db := client.Database(os.Getenv("DBNAME"))

	return db.Collection(name)
}
