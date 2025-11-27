package repository

import (
	"context"
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