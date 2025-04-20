package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetCoinbaseTicker(symbol string) (*Ticker, error) {
	standardMap := map[string]string{
		"BTCUSDT":  "BTC-USDT",
		"ETHUSDT":  "ETH-USDT",
		"BNBUSDT":  "BNB-USDT",
		"XRPUSDT":  "XRP-USDT",
		"DOGEUSDT": "DOGE-USDT",
	}

	symbol = standardMap[symbol]
	url := fmt.Sprintf("https://api.exchange.coinbase.com/products/%s/ticker", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Bid string `json:"bid"`
		Ask string `json:"ask"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	bid, _ := strconv.ParseFloat(result.Bid, 64)
	ask, _ := strconv.ParseFloat(result.Ask, 64)

	return &Ticker{
		Exchange: "Coinbase",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}
