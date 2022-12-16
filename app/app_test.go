package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap/zaptest"

	"github.com/project-code-io/crypto-trading-bot-go/app"
	"github.com/project-code-io/crypto-trading-bot-go/pair"
)

func TestAppStart(t *testing.T) {
	t.Run("app should not get exchange if ran for a short time", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := zaptest.NewLogger(t)

		mockExchange := app.NewmockExchangeClient(ctrl)

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
		mockExchange.EXPECT().GetLastPrice(gomock.Any(), pair.BTCUSD).Times(2).Return("0.00", nil)

		a := app.New(logger, mockExchange)

		ctx, cancel := context.WithCancel(context.Background())

		go a.Start(ctx)
		time.Sleep(time.Second*2 + time.Millisecond*20)
		cancel()
	})
}
