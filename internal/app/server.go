package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/controller"
	"github.com/antnzr/chat-go/internal/app/middleware"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/gin-gonic/gin"
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

// Run starts the HTTP server
// todo: and gRPC server
func (s *Server) Run(ctx context.Context) (err error) {
	var ec = make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	container := buildDeps(s.config, s.pgPool)
	s.httpSrv = NewHttpServer(s.config, *container)

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

type Container struct {
	Controller  controller.Controller
	Middlewares map[string]gin.HandlerFunc
}

func buildDeps(config config.Config, db *pgxpool.Pool) *Container {
	userRepository := repository.NewUserRepository(db)
	tokenRepository := repository.NewTokneRepository(db)
	store := repository.NewStore(userRepository, tokenRepository)

	tokenService := service.NewTokenService(store)
	userService := service.NewUserService(store, tokenService)

	authController := controller.NewAuthController(userService, tokenService)
	userController := controller.NewUserController(userService)

	controller := controller.NewController(authController, userController)
	auth := middleware.Auth(tokenService, userService, config)

	middlewares := make(map[string]gin.HandlerFunc)
	middlewares["auth"] = auth

	return &Container{
		Controller:  *controller,
		Middlewares: middlewares,
	}
}
