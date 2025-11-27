package routes

import (
	"github.com/funcBank_Api/controllers"
	"github.com/gin-gonic/gin"
)

func SeedRoutes(r *gin.Engine, seedController *controllers.SeedController) {
	seedGroup := r.Group("/seed")
	{
		seedGroup.POST("/nav", seedController.SeedNAVData)
	}
}	