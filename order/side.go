package order

// Side represents the side of the trade that the order is placing.
type Side string

const (
	// SideBuy specifies the buy side of an order
	SideBuy Side = "BUY"

	// SideSell specifies the sell side of an order
	SideSell = "SELL"
)
