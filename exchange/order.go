package exchange

import (
	"github.com/project-code-io/crypto-trading-bot-go/order"
	"github.com/project-code-io/crypto-trading-bot-go/trading"
)

// Order represents an order placed on the exchange
type Order struct {
	ID       string
	Pair     trading.Pair
	Side     order.Side
	ClientID string
}
