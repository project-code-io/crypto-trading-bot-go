package exchange

import (
	"context"

	"github.com/project-code-io/crypto-trading-bot-go/order"
	"github.com/project-code-io/crypto-trading-bot-go/trading"
)

// Noop is an exchange that performs no operations on it's functions. This
// type is only used in the scaffolding to help build out the logic.
type Noop struct {
	balances map[trading.Asset]int64
}

func NewNoop() (*Noop, error) {
	const (
		usdAmount = 50
		btcAmount = 0.00001
		ethAmount = 0.05
	)

	return &Noop{
		balances: map[trading.Asset]int64{
			trading.USD: trading.USD.Unit(usdAmount),
			trading.BTC: trading.BTC.Unit(btcAmount),
			trading.ETH: trading.ETH.Unit(ethAmount),
		},
	}, nil
}

// GetLastPrice will return 0.00 for any pair in the noop exchange.
func (e *Noop) GetLastPrice(ctx context.Context, p trading.Pair) (string, error) {
	if p == trading.BTCUSD {
		return "17000.00", nil
	}

	return "5000", nil
}

func (e *Noop) CreateLimitOrder(ctx context.Context, order order.Limit) (Order, error) {
	return Order{}, nil
}

func (e *Noop) CancelOrder(ctx context.Context, orderID string) error {
	return nil
}

func (e *Noop) ListOpenOrders(ctx context.Context) ([]Order, error) {
	return nil, nil
}

func (e *Noop) GetBalance(ctx context.Context, asset trading.Asset) (int64, error) {
	return e.balances[asset], nil
}
