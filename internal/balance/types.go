package balance

type BalanceProvider func(string2 string) *Balance

type Balance struct {
	exchangeName string
	balances     map[string]float64
}
