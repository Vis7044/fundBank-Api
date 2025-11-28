package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/funcBank_Api/models"
	"go.mongodb.org/mongo-driver/mongo"
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
	cursor, err := r.fundCollection.Find(ctx, struct{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var funds []models.SchemeDetail
	for cursor.Next(ctx) {
		var fund models.SchemeDetail
		if err := cursor.Decode(&fund); err != nil {
			return nil, err
		}
		funds = append(funds, fund)
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
