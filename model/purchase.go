package model

type Purchase struct {
	PurchaseID  uint    `json:"id"`
	WagerID     uint    `json:"wager_id"`
	BuyingPrice float64 `json:"buying_price"`
	BoughtAt    int64   `json:"bought_at"`
}
