package main

import (
	"github.com/funcBank_Api/config"
	"github.com/funcBank_Api/controllers"
	"github.com/funcBank_Api/repository"
	"github.com/funcBank_Api/routes"
	"github.com/funcBank_Api/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Connect to database
	config.ConnectDb()
	defer config.DisconnectDatabase()
	db := config.DB.Client().Database("fundBank")

	// Initialize router
	r := gin.Default()

	// regsiter seed routes
	seedRepository := repository.NewSeedRepo(db)
	seedService := services.NewSeedService(seedRepository)
	seedController := controllers.NewSeedController(seedService)
	routes.SeedRoutes(r, seedController)

	// register fund routes
	fundRepository := repository.NewFundRepo(db)
	fundService := services.NewFundService(fundRepository)
	fundController := controllers.NewFundController(fundService)
	routes.FundRoutes(r, fundController)
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
