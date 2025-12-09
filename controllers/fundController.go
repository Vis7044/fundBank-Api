package controllers

import (
	"net/http"

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
	page, err := utils.GetQueryInt64(*ctx, "page", 1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response[string]{Success: false, Data: "Invalid page parameter"})
		return
	}
	limit, err := utils.GetQueryInt64(*ctx, "limit", 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response[string]{Success: false, Data: "Invalid limit parameter"})
		return
	}
	sub_category := ctx.QueryArray("category[]")
	fundhouse := ctx.QueryArray("fundhouse[]")
	// sorting
	sortBy := ctx.DefaultQuery("sortBy", "cagr_1y") // default: nav
	orderStr := ctx.DefaultQuery("order", "desc")

	order := -1
	if orderStr == "asc" {
		order = 1
	}

	funds, err := fc.fundService.GetFunds(ctx, page, limit, sortBy, order, sub_category, fundhouse)
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

func (fc *FundController) SearchFundsByName(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	// pagination
	page, _ := utils.GetQueryInt64(*ctx, "page", 1)
	limit, _ := utils.GetQueryInt64(*ctx, "limit", 10)

	// sorting
	sortBy := ctx.DefaultQuery("sortBy", "nav") // default: nav
	orderStr := ctx.DefaultQuery("order", "desc")

	order := -1
	if orderStr == "asc" {
		order = 1
	}

	results, err := fc.fundService.SearchFundsByName(ctx, query, page, limit, sortBy, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response[[]models.FundScheme]{Success: true, Data: results})
}

func (fc *FundController) SystematicSWPReturnPlan(ctx *gin.Context) {
	type SWPRequest struct {
		SchemeCode          string  `json:"scheme_code" binding:"required"`
		InvestDate          string  `json:"invest_date" binding:"required"`
		SWPDate             int     `json:"swp_date" binding:"required"`
		StartDate           string  `json:"start_date" binding:"required"`
		EndDate             string  `json:"end_date" binding:"required"`
		TotalInvestedAmount float64 `json:"total_invested_amount" binding:"required"`
		WithdrawalAmount    float64 `json:"withdrawal_amount" binding:"required"`
		Interval            string  `json:"interval" binding:"required"` // e.g., "monthly", "quarterly"
	}
	var swpReq SWPRequest
	if err := ctx.ShouldBindJSON(&swpReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"This error": err.Error()})
		return
	}
	//First Call GetFundBySchemeCode to get NAV data for the scheme code within the date range
	fundNavData, err := fc.fundService.GetFundBySchemeCode(ctx, swpReq.SchemeCode, swpReq.InvestDate, swpReq.EndDate)
	fundNavData.Data = utils.ReverseFundNavSlice(fundNavData.Data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// but with weekly, fortnightly,quaterly and monthly intervals
	intervalNav, err := utils.FilterNavByInterval(fundNavData.Data, swpReq.Interval, swpReq.SWPDate, swpReq.StartDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter NAV data by interval"})
		return
	}
	//Then calculate SWP returns based on the NAV data and the SWP parameters provided in the request body
	SWPResponse, err := utils.CalculateSWPResponse(intervalNav, swpReq.TotalInvestedAmount, swpReq.WithdrawalAmount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate SWP returns"})
		return
	}
	//ctx.JSON(http.StatusOK, utils.Response[[]models.SWPResponse]{Success: true, Data: SWPResponse})
	//Print data in tabular format in an excel file and return the file so that user can download it
	filePathForSwp, err := utils.ExportSWPResponseToExcel(SWPResponse, swpReq.SchemeCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate SWP Excel file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":      true,
		"file_path":    filePathForSwp,
		"interval_nav": intervalNav,
		"swp_report":   SWPResponse,
	})

	//also show data in the response body
}
