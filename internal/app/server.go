package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/container"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config  config.Config
	pgPool  *pgxpool.Pool
	httpSrv *HttpServer
	stopFn  sync.Once
}

func NewServer(config config.Config, pgPool *pgxpool.Pool) *Server {
	return &Server{config: config, pgPool: pgPool}
}

// Run the HTTP server
func (s *Server) Run(ctx context.Context) (err error) {
	var ec = make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	container := container.NewContainer(s.config, s.pgPool)
	s.httpSrv = NewHttpServer(s.config, container)

	go func() {
		err := s.httpSrv.Run(ctx)
		if err != nil {
			err = fmt.Errorf("HTTP server error: %w", err)
		}
		ec <- err
	}()

	// Wait for the services to exit.
	var es []string
	for i := 0; i < cap(ec); i++ {
		if err := <-ec; err != nil {
			es = append(es, err.Error())
			// If one of the services returns by a reason other than parent context canceled,
			// try to gracefully shutdown the other services to shutdown everything,
			// with the goal of replacing this service with a new healthy one.
			// NOTE: It might be a slightly better strategy to announce it as unfit for handling traffic,
			// while leaving the program running for debugging.
			if ctx.Err() == nil {
				s.Shutdown(context.Background())
			}
		}
	}
	if len(es) > 0 {
		err = errors.New(strings.Join(es, ", "))
	}
	cancel()
	return err
}

func (s *Server) Shutdown(ctx context.Context) {
	s.stopFn.Do(func() {
		s.httpSrv.Shutdown(ctx)
	})
}
