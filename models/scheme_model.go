package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SchemeDetail struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SchemeCode        string             `bson:"scheme_code" json:"scheme_code"`
	ISINGrowth        string             `bson:"isin_growth" json:"isin_growth"`
	ISINReinvest      string             `bson:"isin_reinvest" json:"isin_reinvest"`
	SchemeName        string             `bson:"scheme_name" json:"scheme_name"`
	ParentName        string             `bson:"parent_name" json:"parent_name"`
	ParentKey         string             `bson:"parent_key" json:"parent_key"`
	FundHouse         string             `bson:"fund_house" json:"fund_house"`
	FundHouseKey      string             `bson:"fund_house_key" json:"fund_house_key"`
	CategoryHeader    string             `bson:"category_header" json:"category_header"`
	Category          string             `bson:"category" json:"category"`
	SubCategory       string             `bson:"sub_category" json:"sub_category"`
	Y1Return          float64            `bson:"1y_return" json:"y1_return"`
	Y3Return          float64            `bson:"3y_return" json:"y3_return"`
	Y5Return          float64            `bson:"5y_return" json:"y5_return"`
	Y1Nav             float64            `bson:"1y_nav" json:"1y_nav"`
	Y3Nav             float64            `bson:"3y_nav" json:"3y_nav"`
	Y5Nav             float64            `bson:"5y_nav" json:"5y_nav"`
	CAGR1             float64            `bson:"cagr_1y" json:"cagr_1y"`
	CAGR3             float64            `bson:"cagr_3y" json:"cagr_3y"`
	CAGR5             float64            `bson:"cagr_5y" json:"cagr_5y"`
	ExpenseRatio      float64            `bson:"expense_ratio" json:"expense_ratio"`
	MinimumInvestment float64            `bson:"minimum_investment" json:"minimum_investment"`
	PlanType          string             `bson:"plan_type" json:"plan_type"`
	OptionType        string             `bson:"option_type" json:"option_type"`
	Frequency         string             `bson:"frequency" json:"frequency"`
	NAV               float64            `bson:"nav" json:"nav"`
	NAVDate           string             `bson:"nav_date" json:"nav_date"`
	DisplayName       string             `bson:"display_name" json:"display_name"`
	AmcImg            string             `bson:"amc_img" json:"amc_img"`
}

type FundResponse struct {
	Meta FundMeta  `json:"meta"`
	Data []FundNav `json:"data"`
}

type FundMeta struct {
	FundHouse           string `json:"fund_house"`
	SchemeType          string `json:"scheme_type"`
	SchemeCategory      string `json:"scheme_category"`
	SchemeCode          int    `json:"scheme_code"`
	SchemeName          string `json:"scheme_name"`
	ISINGrowth          string `json:"isin_growth"`
	ISINDivReinvestment string `json:"isin_div_reinvestment"`
}

type FundNav struct {
	Date string `json:"date"`
	Nav  string `json:"nav"`
}

type FundScheme struct {
	SchemeCode string `bson:"scheme_code" json:"scheme_code"`
	SchemeName string `bson:"scheme_name" json:"scheme_name"`
}

type FundDetail struct {
	SchemeCode     string  `bson:"scheme_code" json:"scheme_code"`
	SchemeName     string  `bson:"scheme_name" json:"scheme_name"`
	ParentName     string  `bson:"parent_name" json:"parent_name"`
	ParentKey      string  `bson:"parent_key" json:"parent_key"`
	FundHouse      string  `bson:"fund_house" json:"fund_house"`
	CategoryHeader string  `bson:"category_header" json:"category_header"`
	Category       string  `bson:"category" json:"category"`
	SubCategory    string  `bson:"sub_category" json:"sub_category"`
	Y1Return          float64            `bson:"1y_return" json:"y1_return"`
	Y3Return          float64            `bson:"3y_return" json:"y3_return"`
	Y5Return          float64            `bson:"5y_return" json:"y5_return"`
	Y1Nav             float64            `bson:"1y_nav" json:"1y_nav"`
	Y3Nav             float64            `bson:"3y_nav" json:"3y_nav"`
	Y5Nav             float64            `bson:"5y_nav" json:"5y_nav"`
	CAGR1             float64            `bson:"cagr_1y" json:"cagr_1y"`
	CAGR3             float64            `bson:"cagr_3y" json:"cagr_3y"`
	CAGR5             float64            `bson:"cagr_5y" json:"cagr_5y"`
	ExpenseRatio   float64 `bson:"expense_ratio" json:"expense_ratio"`
	NAV            float64 `bson:"nav" json:"nav"`
	DisplayName    string  `bson:"display_name" json:"display_name"`
}
