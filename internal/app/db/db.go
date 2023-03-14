package db

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/logger"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

var (
	pgInstance *pgxpool.Pool
	pgOnce     sync.Once
)

func DBPool(ctx context.Context, config config.Config) (*pgxpool.Pool, error) {
	pgOnce.Do(func() {
		db, err := newPGXPool(ctx, config)
		if err != nil {
			logger.Error("unable to create connection pool: %w", zap.Error(err))
			return
		}
		pgInstance = db
	})
	return pgInstance, nil
}

func newPGXPool(ctx context.Context, config config.Config) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	pgxLogLevel, err := logLevelFromEnv(config)
	if err != nil {
		return nil, err
	}

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &PGXStdLogger{},
		LogLevel: pgxLogLevel,
	}

	// pgxpool default max number of connections is the number of CPUs on your machine returned by runtime.NumCPU().
	// This number is very conservative, and you might be able to improve performance for highly concurrent applications
	// by increasing it.
	// conf.MaxConns = int32(runtime.NumCPU() * 2) // maybe later

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}
	return pool, nil
}

// LogLevelFromEnv returns the tracelog.LogLevel from the environment variable PGX_LOG_LEVEL.
// By default this is info (tracelog.LogLevelInfo), which is good for development.
// For deployments, something like tracelog.LogLevelWarn is better choice.
func logLevelFromEnv(config config.Config) (tracelog.LogLevel, error) {
	if level := config.PgLogLevel; level != "" {
		l, err := tracelog.LogLevelFromString(level)
		if err != nil {
			return tracelog.LogLevelDebug, fmt.Errorf("pgx configuration: %w", err)
		}
		return l, nil
	}
	return tracelog.LogLevelInfo, nil
}

// PGXStdLogger prints pgx logs to the logger.
// os.Stderr by default.
type PGXStdLogger struct{}

func (l *PGXStdLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	args := make([]any, 0, len(data)+2) // making space for arguments + level + msg
	args = append(args, level, msg)
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}
	logger.Info(fmt.Sprintln(args...))
}

func PgErrors(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	return fmt.Errorf(`%w
		Code: %v
		Detail: %v
		Hint: %v
		Position: %v
		InternalPosition: %v
		InternalQuery: %v
		Where: %v
		SchemaName: %v
		TableName: %v
		ColumnName: %v
		DataTypeName: %v
		ConstraintName: %v
		File: %v:%v
		Routine: %v`,
		err,
		pgErr.Code,
		pgErr.Detail,
		pgErr.Hint,
		pgErr.Position,
		pgErr.InternalPosition,
		pgErr.InternalQuery,
		pgErr.Where,
		pgErr.SchemaName,
		pgErr.TableName,
		pgErr.ColumnName,
		pgErr.DataTypeName,
		pgErr.ConstraintName,
		pgErr.File, pgErr.Line,
		pgErr.Routine)
}
