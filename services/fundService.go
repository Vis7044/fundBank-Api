package services

import (
	"context"

	"github.com/funcBank_Api/models"
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

func (fs *FundService) GetAllFunds(ctx context.Context) ([]models.SchemeDetail, error) {
	return fs.fundRepo.GetAllFunds(ctx)
}

func (s *FundService) GetFundBySchemeCode(ctx context.Context, schemeCode string, startDate string, endDate string) (*models.FundResponse, error) {
	return s.fundRepo.GetFundBySchemeCode(ctx, schemeCode, startDate, endDate)
}

func (s *FundService) GetFundsByAMC(ctx context.Context, amcName string) ([]models.FundScheme, error) {
	return s.fundRepo.GetFundsByAMC(ctx, amcName)
}

func (fs *FundService) GetFundDetails(ctx context.Context, schemeCode string) (*models.SchemeDetail, error) {
	return fs.fundRepo.GetFundDetails(ctx, schemeCode)
}
