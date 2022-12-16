//go:generate mockgen -source=dependencies.go -destination=./mocks.go -package=app -mock_names ExchangeClient=mockExchangeClient

package app

import (
	"context"

	"github.com/project-code-io/crypto-trading-bot-go/pair"
)

// ExchangeClient represents a type that is able to communicate with an
// exchange.
type ExchangeClient interface {
	GetLastPrice(ctx context.Context, p pair.Pair) (string, error)
}
