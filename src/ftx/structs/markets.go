package structs

import (
	"time"
)

type HistoricalPrice struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	StartTime time.Time `json:"startTime"`
	Volume    float64   `json:"volume"`
}
type HistoricalPrices struct {
	Success bool              `json:"success"`
	Result  []HistoricalPrice `json:"result"`
}

type Trade struct {
	ID          int64     `json:"id"`
	Liquidation bool      `json:"liquidation"`
	Price       float64   `json:"price"`
	Side        string    `json:"side"`
	Size        float64   `json:"size"`
	Time        time.Time `json:"time"`
}

type Trades struct {
	Success bool    `json:"success"`
	Result  []Trade `json:"result"`
}
