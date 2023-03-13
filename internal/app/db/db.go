package db

import (
	"context"
	"sync"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	pgInstance *pgxpool.Pool
	pgOnce     sync.Once
)

func DBPool(config config.Config) (*pgxpool.Pool, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(context.Background(), config.DatabaseURL)
		if err != nil {
			logger.Error("unable to create connection pool: %w", zap.Error(err))
			return
		}
		pgInstance = db
	})
	return pgInstance, nil
}
