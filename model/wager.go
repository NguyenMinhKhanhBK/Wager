package model

import (
	"wager/utils"
)

type Wager struct {
	ID                  uint              `json:"id"`
	TotalWagerValue     uint              `json:"total_wager_value"`
	Odds                uint              `json:"odds"`
	SellingPercentage   uint              `json:"selling_percentage"`
	SellingPrice        float64           `json:"selling_price"`
	CurrentSellingPrice float64           `json:"current_selling_price"`
	PercentageSold      utils.NullUint    `json:"percentage_sold"`
	AmountSold          utils.NullFloat64 `json:"amount_sold"`
	PlaceAt             int64             `json:"place_at"`
}

type CreateWagerRequest struct {
	TotalWagerValue   uint    `json:"total_wager_value" validate:"gt=0"`
	Odds              uint    `json:"odds" validate:"gt=0"`
	SellingPercentage uint    `json:"selling_percentage" validate:"gte=1,lte=100"`
	SellingPrice      float64 `json:"selling_price" validate:"gt=0,monetary-format"`
}

type GetWagerListRequest struct {
	Page  int `validate:"gt=0"`
	Limit int `validate:"gt=0"`
}

type GetWagerListResponse struct {
	Wagers []Wager
}

type BuyWagerRequest struct {
	WagerID     uint    `json:"id"`
	BuyingPrice float64 `json:"buying_price"`
}
