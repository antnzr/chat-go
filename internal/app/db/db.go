package db

import (
	"context"
	"fmt"
	"os"

	"github.com/antnzr/chat-go/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DBPool(config config.Config) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	return db
}
