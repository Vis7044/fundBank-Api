package routes

import (
	"github.com/funcBank_Api/controllers"
	"github.com/gin-gonic/gin"
)

func FundRoutes(r *gin.Engine, fundController *controllers.FundController) {
	fundGroup := r.Group("/funds")
	{
	// 	// fundGroup.GET("/", fundController.GetAllFunds)
	// 	// fundGroup.GET("/:schemeCode", fundController.GetFundBySchemeCode)
	// 	/// get all funds - 200
	// 	// search funds by name
	// 	// get all amc
	// fundGroup.GET("/", fundController.GetAllAMCs)
	// 	// get by scheme code and date range
	fundGroup.GET("/:schemeCode", fundController.GetFundBySchemeCode)
	// 	// get fund by amc
	// 	// get fund deatil
	// 	// get funds by category
	}
}