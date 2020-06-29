package types

// Query endpoints supported by the coins querier
const (
	QueryListCoins = "list"
	QueryGetCoin   = "get"
)

type QueryResCoins []Coin

type QueryResCoin Coin

// implement fmt.Stringer
func (n QueryResCoins) String() string {
	return "to be done"
}

// implement fmt.Stringer
func (n QueryResCoin) String() string {
	return "to be done"
}
