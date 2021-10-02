package data

type InputData struct {
	Symbol           string  `json:"symbol"`
	Price24h        float64 `json:"price_24h"`
	Volume24h       float64 `json:"volume_24h"`
	LastTradePrice float64 `json:"last_trate_price"`
}

type ResultInput struct {
	Result bool `json:"result"`
	PageId uint64 `json:"page,omitempty"`
}

type OutputData struct {
	Price      float64 `json:"price"`
	Volume     float64 `json:"volume"`
	LastTrade float64 `json:"last_trade"`
}
