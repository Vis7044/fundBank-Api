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
