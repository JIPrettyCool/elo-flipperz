package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func GetMongoDBURI() string {
    return os.Getenv("MONGODB_URI")
}