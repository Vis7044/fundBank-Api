package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/funcBank_Api/models"
	"github.com/funcBank_Api/services"
	"github.com/funcBank_Api/utils"
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

func (fc *FundController) GetAllFunds(ctx *gin.Context) {
	// Use a different var name (cxt, mongoCtx, etc.)
	mongoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	funds, err := fc.fundService.GetAllFunds(mongoCtx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response[string]{Success: false, Data: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response[[]models.SchemeDetail]{Success: true, Data: funds})
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

func (fc *FundController) GetAllAMCs(ctx *gin.Context) {
	amcs := []string{
		"All",
		"Aditya Birla Sun Life Mutual Fund",
		"Angel One Mutual Fund",
		"Axis Mutual Fund",
		"Bajaj Finserv Mutual Fund",
		"Bandhan Mutual Fund",
		"Bank of India Mutual Fund",
		"Baroda BNP Paribas Mutual Fund",
		"Canara Robeco Mutual Fund",
		"Capitalmind Mutual Fund",
		"Choice Mutual Fund",
		"DSP Mutual Fund",
		"Edelweiss Mutual Fund",
		"Franklin Templeton Mutual Fund",
		"Groww Mutual Fund",
		"HDFC Mutual Fund",
		"HSBC Mutual Fund",
		"Helios Mutual Fund",
		"ICICI Prudential Mutual Fund",
		"IL&FS Mutual Fund (IDF)",
		"ITI Mutual Fund",
		"Invesco Mutual Fund",
		"JM Financial Mutual Fund",
		"Jio BlackRock Mutual Fund",
		"Kotak Mahindra Mutual Fund",
		"LIC Mutual Fund",
		"Mahindra Manulife Mutual Fund",
		"Mirae Asset Mutual Fund",
		"Motilal Oswal Mutual Fund",
		"NJ Mutual Fund",
		"Navi Mutual Fund",
		"Nippon India Mutual Fund",
		"Old Bridge Mutual Fund",
		"PGIM India Mutual Fund",
		"PPFAS Mutual Fund",
		"Quantum Mutual Fund",
		"SBI Mutual Fund",
		"Samco Mutual Fund",
		"Shriram Mutual Fund",
		"Sundaram Mutual Fund",
		"Tata Mutual Fund",
		"Taurus Mutual Fund",
		"The Wealth Company Mutual Fund",
		"Trust Mutual Fund",
		"UTI Mutual Fund",
		"Unifi Mutual Fund",
		"Union Mutual Fund",
		"WhiteOak Capital Mutual Fund",
		"Zerodha Mutual Fund",
		"quant Mutual Fund",
	}
	ctx.JSON(200, gin.H{"amcs": amcs})
}

func (fc *FundController) GetFundsByAMC(ctx *gin.Context) {
	amcName := ctx.Param("amcName")
	funds, err := fc.fundService.GetFundsByAMC(ctx, amcName)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, funds)
}

func (fc *FundController) GetFundDetails(ctx *gin.Context) {
	var schemeCode string = ctx.Param("schemeCode")
	fund, err := fc.fundService.GetFundDetails(ctx, schemeCode)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response[string]{Success: false, Data: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, utils.Response[*models.FundDetail]{Success: true, Data: fund})

}
