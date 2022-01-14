package db

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// ConnectToDb retrieves db config from .env and tries to conenct to the database.
func ConnectToDb() *mongo.Client {
	if client != nil {
		log.Println("Fetching existing client")
		err := client.Ping(context.Background(), nil)
		if err == nil {
			return client
		}	
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	dbUrl := os.Getenv("DB_URL")

	if username == "" {
		log.Errorln("Missing username in .env")
	}
	if password == "" {
		log.Errorln("Missing password in .env")
	}
	if dbUrl == "" {
		log.Errorln("Missing db url in .env")
	}

	connectionStr := "mongodb+srv://" + username + ":" + password + "@" + dbUrl + "?retryWrites=true&w=majority"
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(connectionStr))
	if err != nil {
		log.Errorln(err)
	}

	// check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Errorln(err)
	}
	log.Println("Connected to mongo client")
	return client
}

// GetMongoDbCollection accepts dbName and collectionname and returns an instance of the specified collection.
func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client := ConnectToDb()

	collection := client.Database(DbName).Collection(CollectionName)
	return collection, nil
}

func DisconnectDB() {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Errorln("FAILED TO CLOSE Mongo Connection")
		log.Errorln(err)
	}

	// TODO optional you can log your closed MongoDB client
	fmt.Println("Connection to MongoDB closed.")

}

