package exchange

import (
	"context"

	"github.com/project-code-io/crypto-trading-bot-go/pair"
)

// Noop is an exchange that performs no operations on it's functions. This
// type is only used in the scaffolding to help build out the logic.
type Noop struct {
}

// GetLastPrice will return 0.00 for any pair in the noop exchange.
func (n *Noop) GetLastPrice(ctx context.Context, p pair.Pair) (string, error) {
	return "0.00", nil
}
