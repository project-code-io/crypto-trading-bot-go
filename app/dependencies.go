//go:generate mockgen -source=dependencies.go -destination=./mocks.go -package=app -mock_names ExchangeClient=mockExchangeClient,IDGenerator=mockIDGenerator

package app

import (
	"context"

	"github.com/project-code-io/crypto-trading-bot-go/exchange"
	"github.com/project-code-io/crypto-trading-bot-go/order"
	"github.com/project-code-io/crypto-trading-bot-go/trading"
)

// ExchangeClient represents a type that is able to communicate with an
// exchange.
type ExchangeClient interface {
	GetLastPrice(ctx context.Context, pair trading.Pair) (string, error)
	CreateLimitOrder(ctx context.Context, order order.Limit) (exchange.Order, error)
	CancelOrders(ctx context.Context, orderIDs ...string) error
	ListOpenOrders(ctx context.Context) ([]exchange.Order, error)
	GetBalance(ctx context.Context, asset trading.Asset) (int64, error)
}

type IDGenerator interface {
	GenerateID(prefix string) string
}
