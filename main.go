package main

import (
	"log"

	"github.com/funcBank_Api/config"
	"github.com/funcBank_Api/controllers"
	"github.com/funcBank_Api/repository"
	"github.com/funcBank_Api/routes"
	"github.com/funcBank_Api/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func startServer(db *mongo.Database) {
	r := gin.Default()

	// seed
	seedRepo := repository.NewSeedRepo(db)
	seedService := services.NewSeedService(seedRepo)
	seedController := controllers.NewSeedController(seedService)
	routes.SeedRoutes(r, seedController)

	// fund
	fundRepo := repository.NewFundRepo(db)
	fundService := services.NewFundService(fundRepo)
	fundController := controllers.NewFundController(fundService)
	routes.FundRoutes(r, fundController)

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	r.Run(":8080")
}

func runCronJobs(fundService *services.FundService) {
	c := cron.New()

	_, err := c.AddFunc("10 21 * * 1-5", func() {
		log.Println("Running daily return calculation...")
		fundService.CalculateReturns() 
		log.Println("Cron Job Completed!")
	})

	if err != nil {
		log.Println("Error scheduling cron:", err)
	}

	c.Start()
}

func main() {
	config.LoadConfig()
	config.ConnectDb()
	db := config.DB.Client().Database("fundBank")

	// Create repository + service ONCE
	fundRepo := repository.NewFundRepo(db)
	fundService := services.NewFundService(fundRepo)

	// Start cron with service
	go runCronJobs(fundService)

	// Start HTTP server
	startServer(db)
}
