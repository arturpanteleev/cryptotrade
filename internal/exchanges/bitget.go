package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetBitgetTicker(symbol string) (*Ticker, error) {
	url := fmt.Sprintf("https://api.bitget.com/api/spot/v1/market/ticker?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			BidPrice string `json:"bestBid"`
			AskPrice string `json:"bestAsk"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no data from Bitget")
	}

	bid, _ := strconv.ParseFloat(result.Data[0].BidPrice, 64)
	ask, _ := strconv.ParseFloat(result.Data[0].AskPrice, 64)

	return &Ticker{
		Exchange: "Bitget",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}
