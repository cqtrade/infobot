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

type LogMessage struct {
	Message string
	Channel string
}

type ReadLogMessage struct {
	Resp chan LogMessage
}

type WriteLogMessage struct {
	Val  LogMessage
	Resp chan bool
}

type ValAt struct {
	Price float64
	At    int64
}
type ReadPriceOp struct {
	Key  string
	Resp chan ValAt
}

type WritePriceOp struct {
	Key  string
	Val  ValAt
	Resp chan bool
}
