package db

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const (
	// Timeout operations after N seconds
	connectTimeout  = 5
	queryTimeout    = 30
	
	// Which instances to read from
	readPreference           = "secondaryPreferred"
	connectionStringTemplate = "mongodb://%s:%s@%s/test?replicaSet=rs0&readpreference=%s&connect=direct&sslInsecure=true&retryWrites=false"
)

// ConnectToDb retrieves db config from .env and tries to conenct to the database.
func ConnectToDb() *mongo.Client {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	dbUrl := os.Getenv("DB_URL")

	if username == "" {
		log.Fatalln("Missing username in .env")
	}
	if password == "" {
		log.Fatalln("Missing password in .env")
	}
	if dbUrl == "" {
		log.Fatalln("Missing db url in .env")
	}

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, dbUrl, readPreference)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}
	
	fmt.Println("Connected to DocumentDB!")

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

