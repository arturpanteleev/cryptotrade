package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetBybitTicker(symbol string) (*Ticker, error) {
	symbol = StandardMap[symbol]

	url := fmt.Sprintf("https://api.bybit.com/v5/market/tickers?category=inverse&symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			List []struct {
				BidPrice string `json:"bid1Price"`
				AskPrice string `json:"ask1Price"`
			} `json:"list"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(result.Result.List) == 0 {
		return nil, fmt.Errorf("no data from Bybit")
	}

	bid, err := strconv.ParseFloat(result.Result.List[0].BidPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid bid price: %w", err)
	}
	ask, err := strconv.ParseFloat(result.Result.List[0].AskPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ask price: %w", err)
	}

	return &Ticker{
		Exchange: "Bybit",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}
