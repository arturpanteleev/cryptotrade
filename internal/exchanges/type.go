package exchange

type Provider func(st string) (*Ticker, error)

type Ticker struct {
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Bid      float64 `json:"bid"`
	Ask      float64 `json:"ask"`
	Time     string  `json:"time"`
}

var StandardMap = map[string]string{
	"BTCUSDT":  "BTCUSDT",
	"ETHUSDT":  "ETHUSDT",
	"BNBUSDT":  "BNBUSDT",
	"XRPUSDT":  "XRPUSDT",
	"DOGEUSDT": "DOGEUSDT",
}
