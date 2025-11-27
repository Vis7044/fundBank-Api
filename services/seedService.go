package services

import (
	"context"

	"github.com/funcBank_Api/repository"
)

type SeedService struct {
	seedRepo *repository.SeedRepo
}

func NewSeedService(seedRepo *repository.SeedRepo) *SeedService {
	return &SeedService{
		seedRepo: seedRepo,
	}
}

func (s *SeedService) SeedNAVData(ctx context.Context, path string) error {
	return s.seedRepo.SeedNAVData(ctx, path)
}