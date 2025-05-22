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
	"BTCUSDT": "BTCUSDT",
	"ETHUSDT": "ETHUSDT",
	"XRPUSDT": "XRPUSDT",
	"TRXUSDT": "TRXUSDT",
	"SOLUSDT": "SOLUSDT",
	"ETHBTC":  "ETHBTC",
	"SOLBTC":  "SOLBTC",
}

var AlterMap = map[string]string{
	"BTCUSDT": "BTC-USDT",
	"ETHUSDT": "ETH-USDT",
	"XRPUSDT": "XRP-USDT",
	"TRXUSDT": "TRX-USDT",
	"SOLUSDT": "SOL-USDT",
	"ETHBTC":  "ETH-BTC",
	"SOLBTC":  "SOL-BTC",
}

var UpBitMap = map[string]string{
	"BTCUSDT": "USDT-BTC",
	"ETHUSDT": "USDT-ETH",
	"XRPUSDT": "USDT-XRP",
	"TRXUSDT": "USDT-TRX",
	"SOLUSDT": "USDT-SOL",
	"ETHBTC":  "BTC-ETH",
	"SOLBTC":  "BTC-SOL",
}
