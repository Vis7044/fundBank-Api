package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SchemeDetail struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SchemeCode     string             `bson:"scheme_code" json:"scheme_code"`
	ISINGrowth     string             `bson:"isin_growth" json:"isin_growth"`
	ISINReinvest   string             `bson:"isin_reinvest" json:"isin_reinvest"`
	SchemeName     string             `bson:"scheme_name" json:"scheme_name"`
	ParentName     string             `bson:"parent_name" json:"parent_name"`
	ParentKey      string             `bson:"parent_key" json:"parent_key"`
	FundHouse      string             `bson:"fund_house" json:"fund_house"`
	FundHouseKey   string             `bson:"fund_house_key" json:"fund_house_key"`
	CategoryHeader string             `bson:"category_header" json:"category_header"`
	CategoryClean  string             `bson:"category_clean" json:"category_clean"`
	PlanType       string             `bson:"plan_type" json:"plan_type"`
	OptionType     string             `bson:"option_type" json:"option_type"`
	Frequency      string             `bson:"frequency" json:"frequency"`
	NAV            float64            `bson:"nav" json:"nav"`
	NAVDate        string             `bson:"nav_date" json:"nav_date"`
	DisplayName    string             `bson:"display_name" json:"display_name"`
	AmcImg         string             `bson:"amc_img" json:"amc_img"`
}
