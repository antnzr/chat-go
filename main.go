package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app"
	"github.com/antnzr/chat-go/internal/app/db"
	"github.com/antnzr/chat-go/internal/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	config, _ := config.LoadConfig(".")
	defer logger.Flush()

	pgPool, err := db.DBPool(context.Background(), config)
	if err != nil {
		logger.Fatality(zap.Error(err))
	}
	defer pgPool.Close()

	srv := app.NewServer(config, pgPool)

	ec := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		ec <- srv.Run(context.Background())
	}()

	// Waits for an internal error that shutdowns the server.
	// Otherwise, wait for a SIGINT or SIGTERM and tries to shutdown the server gracefully.
	// After a shutdown signal, HTTP requests taking longer than the specified grace period are forcibly closed.
	select {
	case err = <-ec:
	case <-ctx.Done():
		haltCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		srv.Shutdown(haltCtx)
		stop()
		err = <-ec
	}

	if err != nil {
		logger.Fatality(zap.Error(err))
	}
}
