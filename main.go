package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Ticker struct {
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Bid      float64 `json:"bid"`
	Ask      float64 `json:"ask"`
	Time     string  `json:"time"`
}

var symbols = []string{"BTCUSDT", "ETHUSDT", "BNBUSDT", "XRPUSDT", "DOGEUSDT"}

func getBinanceTicker(symbol string) (*Ticker, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/bookTicker?symbol=%s", symbol)
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
		Exchange: "Binance",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		Time:     time.Now().Format(time.RFC3339),
	}, nil
}

func getBybitTicker(symbol string) (*Ticker, error) {
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

func getMEXCTicker(symbol string) (*Ticker, error) {
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

func getBitgetTicker(symbol string) (*Ticker, error) {
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

func pricesHandler(w http.ResponseWriter, r *http.Request) {
	var tickers []*Ticker

	for _, symbol := range symbols {
		if t, err := getBinanceTicker(symbol); err == nil {
			tickers = append(tickers, t)
		}
		if t, err := getBybitTicker(symbol); err == nil {
			tickers = append(tickers, t)
		}
		if t, err := getMEXCTicker(symbol); err == nil {
			tickers = append(tickers, t)
		}
		if t, err := getBitgetTicker(symbol); err == nil {
			tickers = append(tickers, t)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickers)
}

func main() {
	http.HandleFunc("/prices", pricesHandler)

	// Статичные файлы (включая index.html)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Server running on http://localhost:88")
	log.Fatal(http.ListenAndServe(":88", nil))
}
