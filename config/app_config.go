package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	MongoURI string
}

var loaded = false
var Cfg *AppConfig

func LoadConfig() {
	if loaded {
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found — using environment variables from Docker or system")
	}

	Cfg = &AppConfig{
		MongoURI: os.Getenv("MONGO_URI"),
	}

	if Cfg.MongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	loaded = true
}
