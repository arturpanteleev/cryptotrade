package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetMEXCTicker(symbol string) (*Ticker, error) {
	url := fmt.Sprintf("https://api.mexc.com/api/v3/ticker/bookTicker?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		BidPrice string `json:"bidPrice"`
		AskPrice string `json:"askPrice"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	bid, _ := strconv.ParseFloat(data.BidPrice, 64)
	ask, _ := strconv.ParseFloat(data.AskPrice, 64)

	return &Ticker{
		Exchange: "MEXC",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}
