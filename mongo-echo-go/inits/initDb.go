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
		_, err = collection.InsertMany(context.TODO(), []interface{}{
			bson.M{
				"status": "active",
				"type":   "1",
				"time":   "1748095225",
				"name":   "sample",
				"value":  "initial data1",
			},
			bson.M{
				"status": "active",
				"type":   "2",
				"time":   "1748008822",
				"name":   "sample",
				"value":  "initial data2",
			},
			bson.M{
				"status": "unactive",
				"type":   "2",
				"time":   "1747922422",
				"name":   "sample",
				"value":  "initial data3",
			},
			bson.M{
				"status": "unactive",
				"type":   "3",
				"time":   "1745330422",
				"name":   "sample",
				"value":  "initial data4",
			},
			bson.M{
				"status": "unactive",
				"type":   "3",
				"time":   "1747490422",
				"name":   "sample",
				"value":  "initial data5",
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted sample data into %s collection\n", collectionName)
	}
}
