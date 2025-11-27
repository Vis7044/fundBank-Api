package repository

import (
	"context"

	"github.com/funcBank_Api/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeedRepo struct {
	fundsCollection *mongo.Collection
}
func NewSeedRepo(db *mongo.Database) *SeedRepo {
	return &SeedRepo{
		fundsCollection: db.Collection("funds"),
	}
}

func (r *SeedRepo) SeedNAVData(ctx context.Context, path string) error {
	records, err := config.ParseNAVAll(path)
	if err != nil {
		return err
	}
	var docs []interface{}
	for _, record := range records {
		docs = append(docs, record)
	}
	_, err = r.fundsCollection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}