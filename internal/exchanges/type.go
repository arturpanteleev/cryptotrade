package exchange

type Provider func(st string) (*Ticker, error)
