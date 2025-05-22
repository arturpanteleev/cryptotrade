package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetOkxTicker(symbol string) (*Ticker, error) {
	symbol = AlterMap[symbol]

	url := fmt.Sprintf("https://www.okx.com/api/v5/market/ticker?instId=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Bid string `json:"bidPx"`
			Ask string `json:"askPx"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no data from OKX")
	}

	bid, _ := strconv.ParseFloat(result.Data[0].Bid, 64)
	ask, _ := strconv.ParseFloat(result.Data[0].Ask, 64)

	return &Ticker{
		Exchange: "OKX",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}
