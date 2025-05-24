package inits

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitCollection(db *mongo.Client, dbName string, collectionName string) {
	collection := db.Database(dbName).Collection(collectionName)

	// Insert sample data if collection is empty
	count, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		_, err = collection.InsertOne(context.TODO(), bson.M{
			"name":  "sample",
			"value": "initial data",
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted sample data into %s collection\n", collectionName)
	}
}
