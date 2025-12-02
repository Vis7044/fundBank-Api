package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	MongoURI string
	Environment string
	Frontend    string
	PORT		string
	DBName		string
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
		Environment: os.Getenv("ENVIRONMENT"),
		Frontend:    os.Getenv("FRONTEND_URL"),
		PORT:		os.Getenv("PORT"),
		DBName:		os.Getenv("DATABASE_NAME"),
	}

	if Cfg.MongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	if Cfg.Environment == "" {
		log.Fatal("ENVIRONMENT is not set")
	}
	loaded = true
}
