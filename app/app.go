package app

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/project-code-io/crypto-trading-bot-go/pair"
)

// App represents the encapsulation of the application state. This struct is
// the main application handler of the trading bot.
type App struct {
	logger   *zap.Logger
	exchange ExchangeClient
	pair     pair.Pair
}

// New acts as the default constructor for the application. Use this method
// to create a new instance of the applcation. This method should
// be called over directly instantiating the App struct as it initializes
// internal attributes so that the application can run as expected.
func New(logger *zap.Logger, exchange ExchangeClient) *App {
	return &App{
		logger:   logger,
		pair:     pair.ETHUSD,
		exchange: exchange,
	}
}

// Start will begin the application with the given context. This method blocks
// based on the context given. In order to stop the application, one should
// cancel the passed in context.
func (a *App) Start(ctx context.Context) {
	a.logger.Info("application starting")

	for {
		select {
		case <-time.After(time.Second):
			price, err := a.exchange.GetLastPrice(ctx, a.pair)
			if err != nil {
				a.logger.Error("failed to get price", zap.Any("pair", a.pair), zap.Error(err))
				break
			}

			a.logger.Info("last price", zap.String("price", price), zap.Any("pair", a.pair))

		case <-ctx.Done():
			a.logger.Info("application shutting down")
			return
		}
	}
}
