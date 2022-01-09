package model

import (
	"database/sql"
	"time"
	"wager/utils"
)

type Wager struct {
	ID                  uint
	TotalWagerValue     uint
	Odds                uint
	SellingPercentage   uint
	SellingPrice        float64
	CurrentSellingPrice float64
	PercentageSold      utils.NullUint
	AmountSold          sql.NullFloat64
	PlaceAt             time.Time
}
