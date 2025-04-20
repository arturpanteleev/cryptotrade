package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetUpbitTicker(symbol string) (*Ticker, error) {
	standardMap := map[string]string{
		"BTCUSDT":  "USDT-BTC",
		"ETHUSDT":  "USDT-ETH",
		"BNBUSDT":  "USDT-BNB",
		"XRPUSDT":  "USDT-XRP",
		"DOGEUSDT": "USDT-DOGE",
	}
	symbol = standardMap[symbol]

	url := fmt.Sprintf("https://api.upbit.com/v1/orderbook?markets=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upbit returned status %s", resp.Status)
	}

	var result []struct {
		OrderbookUnits []struct {
			AskPrice float64 `json:"ask_price"`
			BidPrice float64 `json:"bid_price"`
		} `json:"orderbook_units"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(result) == 0 || len(result[0].OrderbookUnits) == 0 {
		return nil, fmt.Errorf("no data from Upbit orderbook")
	}

	bid := result[0].OrderbookUnits[0].BidPrice
	ask := result[0].OrderbookUnits[0].AskPrice

	return &Ticker{
		Exchange: "Upbit",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}
