package repository

import (
	"context"
	"strconv"
	"strings"

	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/funcBank_Api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FundRepo struct {
	fundCollection *mongo.Collection
}

func NewFundRepo(db *mongo.Database) *FundRepo {
	return &FundRepo{
		fundCollection: db.Collection("funds"),
	}
}

func (r *FundRepo) GetAllFunds(ctx context.Context) ([]models.SchemeDetail, error) {
	opts := options.Find().SetLimit(200)

	cursor, err := r.fundCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var funds []models.SchemeDetail

	// Faster way to decode all results at once
	if err := cursor.All(ctx, &funds); err != nil {
		return nil, err
	}

	return funds, nil
}

func (r *FundRepo) GetFundBySchemeCode(ctx context.Context, schemeCode string, startDate string, endDate string) (*models.FundResponse, error) {

	url := fmt.Sprintf("https://api.mfapi.in/mf/%s?startDate=%s&endDate=%s", schemeCode, startDate, endDate)

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

	return &result, nil
}

func (r *FundRepo) GetFundsByAMC(ctx context.Context, amcName string) ([]models.FundScheme, error) {
	filter := struct {
		Fund_house string `bson:"fund_house"`
	}{
		Fund_house: amcName,
	}

	projection := bson.M{
		"scheme_name": 1,
		"scheme_code": 1,
	}
	cursor, err := r.fundCollection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var funds []models.FundScheme
	for cursor.Next(ctx) {
		var fund models.FundScheme
		if err := cursor.Decode(&fund); err != nil {
			return nil, err
		}
		funds = append(funds, fund)
	}
	return funds, nil
}

func (r *FundRepo) CalculateAndUpdateFundReturns(
    ctx context.Context,
    schemeCode string,
    todayNav *models.FundResponse,
    oneYearNav *models.FundResponse,
    threeYearsNav *models.FundResponse,
    fiveYearsNav *models.FundResponse,
) error {

    returns := make(map[string]float64)

    getNav := func(navResp *models.FundResponse) (float64, bool) {
        if navResp == nil || len(navResp.Data) == 0 {
            return 0, false
        }

        navStr := strings.TrimSpace(navResp.Data[0].Nav)

        navFloat, err := strconv.ParseFloat(navStr, 64)
        if err != nil {
            return 0, false
        }

        return navFloat, true
    }

    today, okToday := getNav(todayNav)
    oneYear, ok1 := getNav(oneYearNav)
    threeYears, ok3 := getNav(threeYearsNav)
    fiveYears, ok5 := getNav(fiveYearsNav)

    // Actual calculations
    if okToday && ok1 && oneYear != 0 {
        returns["1Year"] = ((today - oneYear) / oneYear) * 100
    }

    if okToday && ok3 && threeYears != 0 {
        returns["3Years"] = ((today - threeYears) / threeYears) * 100
    }

    if okToday && ok5 && fiveYears != 0 {
        returns["5Years"] = ((today - fiveYears) / fiveYears) * 100
    }
    update := bson.M{
        "$set": bson.M{
            "1y_return": returns["1Year"],
            "3y_return": returns["3Years"],
            "5y_return": returns["5Years"],
        },
    }

    _, err := r.fundCollection.UpdateOne(
        ctx,
        bson.M{"scheme_code": schemeCode},
        update,
    )

    return err
}

func (fr *FundRepo) GetFundDetails(ctx context.Context, schemeCode string) (*models.SchemeDetail, error) {

	filter := bson.M{"scheme_code": schemeCode}

	result := fr.fundCollection.FindOne(ctx, filter)

	var fund models.SchemeDetail

	if err := result.Decode(&fund); err != nil {
		return nil, err
	}

	return &fund, nil
}
