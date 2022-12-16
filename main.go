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
	l, err := zap.NewProduction()
	if err != nil {
		fmt.Println("error creating logger: %w", err)
	}

	defer func() {
		_ = l.Sync()
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	a := app.New(l, &exchange.Noop{})
	a.Start(ctx)
}
