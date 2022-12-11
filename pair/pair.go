package pair

// Pair represents a currency pairing that can be trading on an exchange.
type Pair string

const (
	// BTCUSD pair represents the BTC/USD pair typically found as a basis
	// in most exchanges.
	BTCUSD Pair = "BTC/USD"

	// ETHUSD pair represents the ETH/USD pair found on exchanges.
	ETHUSD Pair = "ETH/USD"
)
