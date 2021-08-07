package structs

/*

collateralUsed: 2.51464
cost: 25.1464
cumulativeBuySize: 0.009
cumulativeSellSize: 0.001
entryPrice: 3143.3
estimatedLiquidationPrice: 0
future: "ETH-0924"
initialMarginRequirement: 0.1
longOrderSize: 0
maintenanceMarginRequirement: 0.03
netSize: 0.008
openSize: 0.008
realizedPnl: 0.2636
recentAverageOpenPrice: 3116.8
recentBreakEvenPrice: 3113.725
recentPnl: 0.2366
shortOrderSize: 0
side: "buy"
size: 0.008

*/
type Position struct {
	Cost                         float64 `json:"cost"`
	EntryPrice                   float64 `json:"entryPrice"`
	EstimatedLiquidationPrice    float64 `json:"estimatedLiquidationPrice"`
	Future                       string  `json:"future"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
	LongOrderSize                float64 `json:"longOrderSize"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
	NetSize                      float64 `json:"netSize"`
	OpenSize                     float64 `json:"openSize"`
	RealizedPnl                  float64 `json:"realizedPnl"`
	ShortOrderSize               float64 `json:"shortOrderSize"`
	Side                         string  `json:"side"`
	Size                         float64 `json:"size"`
	UnrealizedPnl                float64 `json:"unrealizedPnl"`
	AverageOpenPrice             float64 `json:"recentAverageOpenPrice"`
	BreakEvenPrice               float64 `json:"recentBreakEvenPrice"`
}
type Positions struct {
	Success bool       `json:"success"`
	Result  []Position `json:"result"`
}
