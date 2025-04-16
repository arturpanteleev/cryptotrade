package handlers

import (
	exchange "cryptotrade/internal/exchanges"
	"encoding/json"
	"net/http"
	"sync"
)

var pairs = []string{
	"BTCUSDT",
	"ETHUSDT",
	"BNBUSDT",
	"XRPUSDT",
	"DOGEUSDT",
}

func PricesHandler(w http.ResponseWriter, r *http.Request) {
	type ExchangeData struct {
		Bid float64 `json:"bid"`
		Ask float64 `json:"ask"`
	}

	type SymbolData struct {
		Exchanges      map[string]ExchangeData `json:"exchanges"`
		MinAsk         float64                 `json:"minAsk"`
		MinAskExchange string                  `json:"minAskExchange"`
		MaxBid         float64                 `json:"maxBid"`
		MaxBidExchange string                  `json:"maxBidExchange"`
		Spread         float64                 `json:"spread"`
	}

	result := make(map[string]*SymbolData)

	providers := []exchange.Provider{
		exchange.GetBinanceTicker,
		exchange.GetBitgetTicker,
		exchange.GetBybitTicker,
		exchange.GetMEXCTicker,
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for _, symbol := range pairs {
		for _, f := range providers {
			wg.Add(1)
			provider := f
			s := symbol
			go func() {
				defer wg.Done()
				t, err := provider(s)
				if err != nil {
					return
				}

				mu.Lock()
				defer mu.Unlock()

				sd, ok := result[s]
				if !ok {
					sd = &SymbolData{Exchanges: make(map[string]ExchangeData)}
					result[s] = sd
				}
				sd.Exchanges[t.Exchange] = ExchangeData{Bid: t.Bid, Ask: t.Ask}
			}()
		}
	}

	wg.Wait()

	// Вычисляем minAsk, maxBid, spread
	for _, data := range result {
		minAsk := 1e12
		maxBid := -1.0
		var minAskEx, maxBidEx string

		for ex, val := range data.Exchanges {
			if val.Ask < minAsk {
				minAsk = val.Ask
				minAskEx = ex
			}
			if val.Bid > maxBid {
				maxBid = val.Bid
				maxBidEx = ex
			}
		}

		data.MinAsk = minAsk
		data.MinAskExchange = minAskEx
		data.MaxBid = maxBid
		data.MaxBidExchange = maxBidEx
		if minAsk > 0 {
			data.Spread = ((maxBid - minAsk) / minAsk) * 100
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
