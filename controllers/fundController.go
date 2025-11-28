package controllers

import (
	"github.com/funcBank_Api/services"
	"github.com/gin-gonic/gin"
)

type FundController struct {
	fundService *services.FundService
}
func NewFundController(fundService *services.FundService) *FundController {
	return &FundController{
		fundService: fundService,
	}
}

func (fc *FundController) GetFundBySchemeCode(ctx *gin.Context) {
	schemeCode := ctx.Param("schemeCode")

	type DateRange struct {
		StartDate string `form:"startDate" binding:"required"`
		EndDate   string `form:"endDate" binding:"required"`
	}

	var dr DateRange

	if err := ctx.ShouldBindQuery(&dr); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid date range parameters"})
		return
	}

	fund, err := fc.fundService.GetFundBySchemeCode(ctx, schemeCode, dr.StartDate, dr.EndDate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, fund)
}
