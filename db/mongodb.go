package db

import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "elo-flipperz/config"
)

var (
    Client *mongo.Client
)

func Connect() {
    MONGODB_URI := config.GetMongoDBURI()
    if MONGODB_URI == "" {
        log.Fatal("MONGODB_URI environment variable not set")
    }

    clientOptions := options.Client().ApplyURI(MONGODB_URI)
    var err error
    Client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    log.Println("Connected MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
    return Client.Database("coinflip").Collection(collectionName)
}