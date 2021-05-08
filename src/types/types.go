package types

type JSONMessageBody struct {
	Ticker   string  `json:"ticker"`
	Exchange string  `json:"exchange"`
	Signal   float64 `json:"signal"`
	Ch24h    float64 `json:"change_24h"`
	Ch7      float64 `json:"change_7d"`
}
