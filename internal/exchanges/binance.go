package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetBinanceTicker(symbol string) (*Ticker, error) {
	symbol = StandardMap[symbol]

	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/bookTicker?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		BidPrice string `json:"bidPrice"`
		AskPrice string `json:"askPrice"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	bid, err := strconv.ParseFloat(result.BidPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid bid: %w", err)
	}
	ask, err := strconv.ParseFloat(result.AskPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ask: %w", err)
	}

	return &Ticker{
		Exchange: "Binance",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}
