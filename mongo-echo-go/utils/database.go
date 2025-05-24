package utils

import (
	"context"
	"fmt"
	"log"

	"mongo-echo-go/config"
	"mongo-echo-go/inits"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client
var DBName string
var CollectionName string

func InitDB() {
	// Get config
	mongoConfig := config.GetMongoConfig()
	DBName = mongoConfig.DBName
	CollectionName = mongoConfig.Collection

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoConfig.URI())

	// Connect to MongoDB
	var err error
	DB, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = DB.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Initialize collection with sample data
	inits.InitCollection(DB, DBName, CollectionName)
}

func GetDB() *mongo.Client {
	return DB
}

func GetDBName() string {
	return DBName
}

func GetCollectionName() string {
	return CollectionName
}
