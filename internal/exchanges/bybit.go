package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetBybitTicker(symbol string) (*Ticker, error) {
	url := fmt.Sprintf("https://api.bybit.com/v2/public/tickers?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result []struct {
			BidPrice string `json:"bid_price"`
			AskPrice string `json:"ask_price"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Result) == 0 {
		return nil, fmt.Errorf("no data from Bybit")
	}

	bid, _ := strconv.ParseFloat(result.Result[0].BidPrice, 64)
	ask, _ := strconv.ParseFloat(result.Result[0].AskPrice, 64)

	return &Ticker{
		Exchange: "Bybit",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}
