package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap/zaptest"

	"github.com/project-code-io/crypto-trading-bot-go/app"
	"github.com/project-code-io/crypto-trading-bot-go/exchange"
	"github.com/project-code-io/crypto-trading-bot-go/order"
	"github.com/project-code-io/crypto-trading-bot-go/trading"
)

func TestAppStart(t *testing.T) {
	t.Run("app should not get exchange if ran for a short time", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := zaptest.NewLogger(t)

		mockExchange := app.NewmockExchangeClient(ctrl)
		mockExchange.EXPECT().ListOrders(gomock.Any()).Return([]exchange.Order{}, nil)

		a := app.New(logger, mockExchange)

		ctx, cancel := context.WithCancel(context.Background())

		go a.Start(ctx)
		time.Sleep(time.Millisecond * 500)
		cancel()
	})

	t.Run("app should call get exchange once per second", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := zaptest.NewLogger(t)

		mockExchange := app.NewmockExchangeClient(ctrl)
		mockExchange.EXPECT().ListOrders(gomock.Any()).Return([]exchange.Order{}, nil)
		mockExchange.EXPECT().GetLastPrice(gomock.Any(), trading.BTCUSD).Times(2).Return("1000.00", nil)
		mockExchange.EXPECT().GetBalance(gomock.Any(), trading.USD).Times(2).Return(int64(5000), nil)

		mockExchange.EXPECT().CreateLimitOrder(gomock.Any(), order.Limit{
			ClientID: "foobar",
			Pair:     trading.BTCUSD,
			Side:     order.SideBuy,
			BaseSize: "0.01",
			Price:    "500",
			PostOnly: true,
		}).Times(2).Return(exchange.Order{}, nil)

		idGen := app.NewmockIDGenerator(ctrl)
		idGen.EXPECT().GenerateID("go-trading-bot").Times(2).Return("foobar")

		a := app.New(logger, mockExchange, app.WithIDGenerator(idGen))

		ctx, cancel := context.WithCancel(context.Background())

		go a.Start(ctx)

		time.Sleep(time.Second*2 + time.Millisecond*20)
		cancel()
	})
}
