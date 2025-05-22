package handlers

import (
	exchange "cryptotrade/internal/exchanges"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

const comission = 0.001

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

func PricesHandler(w http.ResponseWriter, r *http.Request) {

	result := make(map[string]*SymbolData)

	providers := []exchange.Provider{
		exchange.GetBinanceTicker,
		exchange.GetBitgetTicker,
		exchange.GetBybitTicker,
		exchange.GetMEXCTicker,
		exchange.GetOkxTicker,
		exchange.GetUpbitTicker,
		exchange.GetCoinbaseTicker,
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for _, symbol := range exchange.StandardMap {
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

				if t.Bid == 0 && t.Ask == 0 {
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
			//todo cсделать отправку в телегу
			if data.Spread > 0.2 {

			}

		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Функция форматирования сообщения
func formatTelegramMessage(symbol string, data *SymbolData) string {
	message := fmt.Sprintf(
		"📈 Символ: %s\n"+
			"Минимальный Ask: %.4f (%s)\n"+
			"Максимальный Bid: %.4f (%s)\n"+
			"➖ Спред: %.4f",
		symbol,
		data.MinAsk, data.MinAskExchange,
		data.MaxBid, data.MaxBidExchange,
		data.Spread,
	)

	return message
}
