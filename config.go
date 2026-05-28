package main

import (
	"os"
	"path/filepath"

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
	exe, err := os.Executable()
	if err != nil {
		return Config{}, err
	}
	if err := godotenv.Load(filepath.Join(filepath.Dir(exe), ".env")); err != nil {
		return Config{}, err
	}

	return Config{
		MongoUri:       os.Getenv("MONGO_URI"),
		DatabaseName:   os.Getenv("DATABASE_NAME"),
		CollectionName: os.Getenv("COLLECTION_NAME"),
	}, nil
}
