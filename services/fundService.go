package services

import (
	"github.com/funcBank_Api/repository"
)

type FundService struct {
	fundRepo *repository.FundRepo
}
func NewFundService(fundRepo *repository.FundRepo) *FundService {
	return &FundService{
		fundRepo: fundRepo,
	}
}