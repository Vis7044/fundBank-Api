package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/funcBank_Api/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SaveJSON(filename string, data interface{}) error {
	// Pretty printed JSON
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, b, 0644)
}


func main() {
	// Load environment variables
	config.LoadConfig()

	// Connect to database
	config.ConnectDb()
	defer config.DisconnectDatabase()

	// Initialize router
	r := gin.Default()

	records, err := config.ParseNAVAll("funds.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Take first 10 safely
	limit := 10
	if len(records) < 10 {
		limit = len(records)
	}
	firstTen := records[:limit]

	// Print to console
	for _, rec := range firstTen {
		fmt.Printf("%+v\n", rec)
	}

	// Save to JSON file
	err = SaveJSON("first10.json", firstTen)
	if err != nil {
		log.Fatalf("failed to save json: %v", err)
	}

	fmt.Println("Saved first 10 records to first10.json")

	// Allow CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// Run server
	r.Run(":8080")
}
