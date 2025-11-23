package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Scheme struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SchemeCode    int                `bson:"scheme_code" json:"scheme_code"`
	FundHouse     string             `bson:"fund_house" json:"fund_house"`
	Category      string             `bson:"category" json:"category"`
	SchemeName    string             `bson:"scheme_name" json:"scheme_name"`
	PlanType      string             `bson:"plan_type" json:"plan_type"`
	OptionType    string             `bson:"option_type" json:"option_type"`
	IsActive      bool               `bson:"is_active" json:"is_active"`
	LatestNAV     float64            `bson:"latest_nav" json:"latest_nav"`
	LatestNAVDate string             `bson:"latest_nav_date" json:"latest_nav_date"`
	LogoURL       string             `bson:"logo_url" json:"logo_url"`
}
