package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ParentFund struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ParentFundKey     string             `bson:"parent_fund_key" json:"parent_fund_key"`
	FundHouse         string             `bson:"fund_house" json:"fund_house"`
	Category          string             `bson:"category" json:"category"`
	ActiveSchemeCodes []int              `bson:"active_scheme_codes" json:"active_scheme_codes"`
}

