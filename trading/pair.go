package trading

// Pair represents an asset pairing that can be trading on an exchange.
type Pair struct {
	Base  Asset
	Quote Asset
}

var (
	// BTCUSD pair represents the BTC/USD pair typically found as a basis
	// in most exchanges.
	BTCUSD = Pair{
		Base:  BTC,
		Quote: USD,
	}

	// ETHUSD pair represents the ETH/USD pair found on exchanges.
	ETHUSD = Pair{
		Base:  ETH,
		Quote: USD,
	}
)
