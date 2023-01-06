package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/project-code-io/crypto-trading-bot-go/app"
	"github.com/project-code-io/crypto-trading-bot-go/exchange"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("error creating logger:", err)
		return
	}

	defer func() {
		_ = logger.Sync()
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	e := exchange.NewBinance(exchange.BinanceDomainUS)
	a := app.New(logger, e)
	a.Start(ctx)
}
