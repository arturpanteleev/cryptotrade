package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Структуры для разбора ответов с бирж
type binanceTicker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type bybitV5Response struct {
	Result struct {
		List []struct {
			Symbol    string `json:"symbol"`
			LastPrice string `json:"lastPrice"`
		} `json:"list"`
	} `json:"result"`
}

type okxResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstId string `json:"instId"`
		Last   string `json:"last"`
	} `json:"data"`
}

type kucoinResponse struct {
	Code string `json:"code"`
	Data struct {
		Price string `json:"price"`
	} `json:"data"`
}

type coingeckoResponse map[string]map[string]float64

type coincapResponse struct {
	Data struct {
		PriceUsd string `json:"priceUsd"`
	} `json:"data"`
}

type coinloreResponse []struct {
	PriceUsd string `json:"price_usd"`
}

type cryptocompareResponse map[string]float64

// Функции для получения цен с бирж
func getBinance(symbol string) string {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка Binance:", err)
		return "error"
	}
	defer resp.Body.Close()

	var t binanceTicker
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		log.Println("Ошибка декодирования Binance:", err)
		return "error"
	}
	return t.Price
}

func getBybit(symbol string) string {
	url := fmt.Sprintf("https://api.bybit.com/v5/market/tickers?category=linear&symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка Bybit:", err)
		return "error"
	}
	defer resp.Body.Close()

	var res bybitV5Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Ошибка декодирования Bybit:", err)
		return "error"
	}
	if len(res.Result.List) > 0 {
		return res.Result.List[0].LastPrice
	}
	return "not found"
}

func getOKX(symbol string) string {
	url := fmt.Sprintf("https://www.okx.com/api/v5/market/tickers?instType=SPOT")
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка OKX:", err)
		return "error"
	}
	defer resp.Body.Close()

	var res okxResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Ошибка декодирования OKX:", err)
		return "error"
	}
	for _, item := range res.Data {
		if item.InstId == symbol {
			return item.Last
		}
	}
	return "not found"
}

func getKuCoin(symbol string) string {
	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка KuCoin:", err)
		return "error"
	}
	defer resp.Body.Close()

	var res kucoinResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Ошибка декодирования KuCoin:", err)
		return "error"
	}
	return res.Data.Price
}

func getCoinGecko(coinID string) string {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coinID)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка CoinGecko:", err)
		return "error"
	}
	defer resp.Body.Close()

	var res coingeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Ошибка декодирования CoinGecko:", err)
		return "error"
	}
	if price, ok := res[coinID]["usd"]; ok {
		return fmt.Sprintf("%.2f", price)
	}
	return "not found"
}

func getCoinCap(assetID string) string {
	url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s", assetID)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка CoinCap:", err)
		return "error"
	}
	defer resp.Body.Close()

	var res coincapResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Ошибка декодирования CoinCap:", err)
		return "error"
	}
	return res.Data.PriceUsd
}

func getCoinLore(id string) string {
	url := fmt.Sprintf("https://api.coinlore.net/api/ticker/?id=%s", id)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка CoinLore:", err)
		return "error"
	}
	defer resp.Body.Close()

	var res coinloreResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Ошибка декодирования CoinLore:", err)
		return "error"
	}
	if len(res) > 0 {
		return res[0].PriceUsd
	}
	return "not found"
}

func getCryptoCompare(symbol string) string {
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/price?fsym=%s&tsyms=USD", symbol)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка CryptoCompare:", err)
		return "error"
	}
	defer resp.Body.Close()

	var res cryptocompareResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Ошибка декодирования CryptoCompare:", err)
		return "error"
	}
	if price, ok := res["USD"]; ok {
		return fmt.Sprintf("%.2f", price)
	}
	return "not found"
}

// HTTP-обработчик
func handler(w http.ResponseWriter, r *http.Request) {
	rates := map[string]string{
		"Binance BTCUSDT":      getBinance("BTCUSDT"),
		"Bybit BTCUSDT":        getBybit("BTCUSDT"),
		"OKX BTC-USDT":         getOKX("BTC-USDT"),
		"KuCoin BTC-USDT":      getKuCoin("BTC-USDT"),
		"CoinGecko bitcoin":    getCoinGecko("bitcoin"),
		"CoinCap bitcoin":      getCoinCap("bitcoin"),
		"CoinLore BTC (id:90)": getCoinLore("90"),
		"CryptoCompare BTC":    getCryptoCompare("BTC"),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rates); err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	port := 88
	log.Printf("Сервер запущен на порту %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
