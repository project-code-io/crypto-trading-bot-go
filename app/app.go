package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/project-code-io/crypto-trading-bot-go/generator"
	"github.com/project-code-io/crypto-trading-bot-go/order"
	"github.com/project-code-io/crypto-trading-bot-go/trading"
)

// App represents the encapsulation of the application state. This struct is
// the main application handler of the trading bot.
type App struct {
	logger      *zap.Logger
	exchange    ExchangeClient
	pair        trading.Pair
	prefix      string
	idGenerator IDGenerator
}

// New acts as the default constructor for the application. Use this method
// to create a new instance of the applcation. This method should
// be called over directly instantiating the App struct as it initializes
// internal attributes so that the application can run as expected.
func New(logger *zap.Logger, exchange ExchangeClient, opts ...Option) *App {
	app := &App{
		logger:      logger,
		pair:        trading.BTCUSD,
		exchange:    exchange,
		prefix:      "go-trading-bot",
		idGenerator: &generator.RandomUUIDGenerator{},
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// Start will begin the application with the given context. This method blocks
// based on the context given. In order to stop the application, one should
// cancel the passed in context.
func (a *App) Start(ctx context.Context) {
	a.logger.Info("application starting")

	if err := a.clearOldOrders(ctx); err != nil {
		a.logger.Error("could not clear old olders", zap.Error(err))
		return
	}

	for {
		select {
		case <-time.After(time.Second):
			price, err := a.exchange.GetLastPrice(ctx, a.pair)
			if err != nil {
				a.logger.Error("failed to get price", zap.Any("pair", a.pair), zap.Error(err))
				break
			}

			a.logger.Info("last price", zap.String("price", price), zap.Any("pair", a.pair))

			if err := a.createAndClearOrder(ctx, price); err != nil {
				a.logger.Error("failed to create and clear order, exiting early", zap.Error(err))
				return
			}
		case <-ctx.Done():
			a.logger.Info("application shutting down")
			return
		}
	}
}

func (a *App) clearOldOrders(ctx context.Context) error {
	a.logger.Info("clearing old orders")

	orders, err := a.exchange.ListOpenOrders(ctx)
	if err != nil {
		return fmt.Errorf("list all orders: %w", err)
	}

	orderIDs := make([]string, 0)

	for _, order := range orders {
		if !strings.HasPrefix(order.ClientID, a.prefix) {
			continue
		}

		orderIDs = append(orderIDs, order.ID)
	}

	if err := a.exchange.CancelOrders(ctx, orderIDs...); err != nil {
		return fmt.Errorf("cancel orders: %w", err)
	}

	a.logger.Info("orders cleared")

	return nil
}

func (a *App) createAndClearOrder(ctx context.Context, price string) error {
	const fundUse = float64(0.1)

	a.logger.Info("loading balance")

	balance, err := a.exchange.GetBalance(ctx, a.pair.Quote)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}

	quoteAmount := balance / int64(1/fundUse)

	quotePrice, err := a.pair.Quote.UnitStr(price)
	if err != nil {
		return fmt.Errorf("quote price: %w", err)
	}

	// Try to buy at half the current price.
	const priceDivisor = 2
	desiredPrice := quotePrice / priceDivisor

	baseSize := a.pair.Base.Unit(float64(quoteAmount) / float64(desiredPrice))

	fmt.Println(quoteAmount, "/", desiredPrice, "=", baseSize)

	o := order.Limit{
		ClientID: a.idGenerator.GenerateID(a.prefix),
		Pair:     a.pair,
		Side:     order.SideBuy,
		BaseSize: a.pair.Base.Format(baseSize),
		Price:    a.pair.Quote.Format(desiredPrice),
		PostOnly: true,
	}

	a.logger.Info("creating order", zap.Any("order", o))

	eOrder, err := a.exchange.CreateLimitOrder(ctx, o)
	if err != nil {
		return fmt.Errorf("create limit order: %w", err)
	}

	a.logger.Info("order created", zap.Any("exchange_order", eOrder))

	const waitTime = time.Millisecond * 200

	select {
	case <-time.After(waitTime):
		err := a.exchange.CancelOrders(ctx, eOrder.ID)
		if err != nil {
			return fmt.Errorf("cancel order: %w", err)
		}
	case <-ctx.Done():
		break
	}

	return nil
}
