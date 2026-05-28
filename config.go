package main

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the connection parameters for MongoDB.
type Config struct {
	MongoUri       string
	DatabaseName   string
	CollectionName string
}

// NewConfig loads configuration from a .env file and returns a populated Config.
func NewConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	return Config{
		MongoUri:       os.Getenv("MONGO_URI"),
		DatabaseName:   os.Getenv("DATABASE_NAME"),
		CollectionName: os.Getenv("COLLECTION_NAME"),
	}, nil
}
