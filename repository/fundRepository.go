package repository

import (
	"context"

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

