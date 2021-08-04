package types

type JSONMessageBody struct {
	Signal   float64 `json:"sig,omitempty"`
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
