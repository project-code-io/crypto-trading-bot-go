package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
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

	if err = godotenv.Load(); err != nil {
		logger.Error("failed to load dotenv", zap.Error(err))
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	exc, err := exchange.NewNoop()
	if err != nil {
		logger.Error("failed to load exchange", zap.Error(err))
		return
	}

	a := app.New(logger, exc)
	a.Start(ctx)
}
