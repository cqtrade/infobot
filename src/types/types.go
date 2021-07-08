package types

/*
{
	"signal":{{plot("Signal")}},
	"ticker":"{{ticker}}",
	"atrtp":{{plot("AtrTp")}},
	"atrsl":{{plot("AtrSl")}},
	"risk":{{plot("Risk")}}",
	"exchange":"{{exchange}}"
}
*/
type JSONMessageBody struct {
	Signal   float64 `json:"signal,omitempty"`
	Ticker   string  `json:"ticker,omitempty"`
	AtrTP    float64 `json:"atrtp,omitempty"`
	AtrSL    float64 `json:"atrsl,omitempty"`
	Risk     float64 `json:"risk,omitempty"`
	Exchange string  `json:"exchange,omitempty"`
	Text     string  `json:"text,omitempty"`
	Sub      string  `json:"sub,omitempty"`
}

type NotificationBody struct {
	Content string `json:"content"`
}
