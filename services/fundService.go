package services

import (
	"context"

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

func (s *FundService) GetFundBySchemeCode(ctx context.Context, schemeCode string, startDate string, endDate string) (interface{}, error) {
	return s.fundRepo.GetFundBySchemeCode(ctx, schemeCode, startDate, endDate)
}