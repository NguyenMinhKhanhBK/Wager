package model

import "time"

type Purchase struct {
	ID          uint
	WagerID     uint
	Wager       Wager
	BuyingPrice float64
	BoughtAt    time.Time
}
