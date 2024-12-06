package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

const (
	databasePrefix   = "mongo_rs_checker"
	collectionPrefix = "mongo_rs_checker"
)

func generateSuffix(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("connectToMongoDB() successful")
	return client, nil
}

func queryMongoDB(client *mongo.Client, databaseName, collectionName string) ([]interface{}, error) {
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var results []interface{}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func insertDocument(client *mongo.Client, databaseName, collectionName string, document map[string]interface{}) (interface{}, error) {
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	insertResult, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)
	return insertResult.InsertedID, nil
}

func findDocumentByID(client *mongo.Client, databaseName, collectionName string, id interface{}) (map[string]interface{}, error) {
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	var result map[string]interface{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	log.Printf("Found document: %v\n", result)
	return result, nil
}

func dropCollection(client *mongo.Client, databaseName, collectionName string) error {
	collection := client.Database(databaseName).Collection(collectionName)
	err := collection.Drop(context.TODO())
	if err != nil {
		return err
	}

	log.Printf("Collection %s dropped successfully\n", collectionName)
	return nil
}

func main() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://user:passw1@localhost:27017"
	}

	client, err := connectToMongoDB(uri)
	if err != nil {
		log.Println("connectToMongoDB error:", err)
		log.Fatal(err)
		return
	}
	defer client.Disconnect(context.TODO())

	databaseName := os.Getenv("MONGODB_DATABASE_NAME")
	if databaseName == "" {
		// Generate a random database name if not set
		randomSuffix, err := generateSuffix(4)
		if err != nil {
			fmt.Println("Error generateSuffix:", err)
			return
		}
		databaseName = databasePrefix + "_" + randomSuffix
	}

	// Always generate a random collection name
	randomSuffix, err := generateSuffix(4)
	if err != nil {
		fmt.Println("Error generateSuffix:", err)
		return
	}

	collectionName := collectionPrefix + "_" + randomSuffix

	log.Printf("Using database: %s", databaseName)
	log.Printf("Using collection: %s", collectionName)

	document := map[string]interface{}{
		"name": "Bi Ba",
		"age":  30,
		"city": "NY",
	}

	insertedID, err := insertDocument(client, databaseName, collectionName, document)
	if err != nil {
		log.Println("insertDocument error:", err)
		log.Fatal(err)
		return
	}

	_, err = findDocumentByID(client, databaseName, collectionName, insertedID)
	if err != nil {
		log.Println("findDocumentByID error:", err)
		log.Fatal(err)
		return
	}

	results, err := queryMongoDB(client, databaseName, collectionName)
	if err != nil {
		log.Println("queryMongoDB error:", err)
		log.Fatal(err)
		return
	}

	for _, result := range results {
		fmt.Println(result)
	}

	err = dropCollection(client, databaseName, collectionName)
	if err != nil {
		log.Println("dropCollection error:", err)
		log.Fatal(err)
		return
	}
}
