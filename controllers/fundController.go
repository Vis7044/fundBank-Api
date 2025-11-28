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
