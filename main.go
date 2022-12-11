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

	a := app.New(logger, &exchange.Noop{})
	a.Start(ctx)
}
