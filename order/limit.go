package order

import (
	"time"

	"github.com/project-code-io/crypto-trading-bot-go/trading"
)

// Limit represents an order that a limit type
type Limit struct {
	ClientID string
	Pair     trading.Pair
	Side     Side
	BaseSize string
	Price    string
	PostOnly bool
	Expires  *time.Time
}
