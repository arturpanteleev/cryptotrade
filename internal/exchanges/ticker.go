package exchange

type Ticker struct {
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Bid      float64 `json:"bid"`
	Ask      float64 `json:"ask"`
	Time     string  `json:"time"`
}
