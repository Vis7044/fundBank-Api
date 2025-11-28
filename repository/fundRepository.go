package repository

import (
	"context"

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
