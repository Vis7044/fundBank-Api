package controllers

import (
	"github.com/funcBank_Api/services"
)

type FundController struct {
	fundService *services.FundService
}
func NewFundController(fundService *services.FundService) *FundController {
	return &FundController{
		fundService: fundService,
	}
}
