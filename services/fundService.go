package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

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

func (fs *FundService) GetFunds(ctx context.Context, page, limit int64, sub_category string) ([]models.SchemeDetail, error) {
	return fs.fundRepo.GetFunds(ctx, page, limit, sub_category)
}

func (s *FundService) GetFundBySchemeCode(ctx context.Context, schemeCode string, startDate string, endDate string) (*models.FundResponse, error) {
	return s.fundRepo.GetFundBySchemeCode(ctx, schemeCode, startDate, endDate)
}

func (s *FundService) GetFundsByAMC(ctx context.Context, amcName string) ([]models.FundScheme, error) {
	return s.fundRepo.GetFundsByAMC(ctx, amcName)
}

// CalculateReturns processes all funds in batches with controlled concurrency
func (s *FundService) CalculateReturns() error {
	ctx := context.Background()

	funds, err := s.fundRepo.GetAllFunds(ctx)
	if err != nil {
		return err
	}

	batchSize := 500
	workerLimit := 30
	totalFunds := len(funds)

	for start := 0; start < totalFunds; start += batchSize {
		end := start + batchSize
		if end > totalFunds {
			end = totalFunds
		}

		batch := funds[start:end]

		if err := s.processBatch(ctx, batch, workerLimit); err != nil {
			return err
		}
	}

	return nil
}

// processBatch handles a batch of funds with controlled concurrency
// using a semaphore pattern
// to limit the number of concurrent goroutines.
// Explanation : A buffered channel (sem) is used as a semaphore to limit the number of concurrent workers.
// Before starting a goroutine, we send a value into the channel.
// When a goroutine finishes, it reads from the channel to free up a slot.
// This ensures that at most 'workerLimit' goroutines are running simultaneously.
// The WaitGroup (wg) is used to wait for all goroutines in the batch to complete before returning.
func (s *FundService) processBatch(ctx context.Context, batch []models.SchemeDetail, workerLimit int) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, workerLimit)

	for _, fund := range batch {

		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			if err := s.CalculateFundReturns(ctx, fund.SchemeCode); err != nil {
				fmt.Println("Error calculating returns for", fund.SchemeCode, ":", err)
			}
		}()
	}

	wg.Wait()
	return nil
}

func CheckWeekend(date time.Time) time.Time {
	switch date.Weekday() {
	case time.Saturday:
		date = date.AddDate(0, 0, -1)
	case time.Sunday:
		date = date.AddDate(0, 0, -2)
	}
	return date
}

func LatestNav(shcemeCode string) (*models.FundResponse, error) {
	url := fmt.Sprintf("https://api.mfapi.in/mf/%s/latest", shcemeCode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result models.FundResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no data available")
	}
	return &result, nil
}

func (s *FundService) CalculateFundReturns(ctx context.Context, schemeCode string) error {
	today := time.Now()

	oneYearAgoStr := CheckWeekend(today.AddDate(-1, 0, 0)).Format("2006-01-02")
	threeYearsAgoStr := CheckWeekend(today.AddDate(-3, 0, 0)).Format("2006-01-02")
	fiveYearsAgoStr := CheckWeekend(today.AddDate(-5, 0, 0)).Format("2006-01-02")

	// Fetch NAVs
	todayNav, err := LatestNav(schemeCode)
	if err != nil {
		return err
	}

	oneYearNav, err := s.fundRepo.GetFundBySchemeCode(ctx, schemeCode, oneYearAgoStr, oneYearAgoStr)
	if err != nil {
		return err
	}

	threeYearsNav, err := s.fundRepo.GetFundBySchemeCode(ctx, schemeCode, threeYearsAgoStr, threeYearsAgoStr)
	if err != nil {
		return err
	}

	fiveYearsNav, err := s.fundRepo.GetFundBySchemeCode(ctx, schemeCode, fiveYearsAgoStr, fiveYearsAgoStr)
	if err != nil {
		return err
	}

	// Calculate and update DB
	return s.fundRepo.CalculateAndUpdateFundReturns(
		context.Background(),
		schemeCode,
		todayNav,
		oneYearNav,
		threeYearsNav,
		fiveYearsNav,
	)
}
func (fs *FundService) GetFundDetails(ctx context.Context, schemeCode string) (*models.FundDetail, error) {
	return fs.fundRepo.GetFundDetails(ctx, schemeCode)
}
