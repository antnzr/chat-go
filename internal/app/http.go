package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/middleware"
	"github.com/antnzr/chat-go/internal/pkg/logger"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	config    config.Config
	container Container
	http      *http.Server
}

func NewHttpServer(
	config config.Config,
	container Container,
) *HttpServer {
	return &HttpServer{
		config:    config,
		container: container,
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	engine := s.setup()

	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.config.Port),
		Handler: engine,
	}

	if err := s.http.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) {
	logger.Info("shutting down HTTP server")
	if s.http != nil {
		if err := s.http.Shutdown(ctx); err != nil {
			logger.Error("failed graceful shutdown of HTTP server")
		}
	}
}

func (s *HttpServer) setup() *gin.Engine {
	gin.SetMode(s.config.GinMode)

	engine := gin.New()
	engine.SetTrustedProxies(nil)
	engine.Use(middleware.ErrorHandler())

	engine.Use(ginzap.Ginzap(logger.GetLogger(), time.RFC3339Nano, true))
	engine.Use(ginzap.RecoveryWithZap(logger.GetLogger(), true))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = s.config.Origin
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))

	s.setupRoutes(engine)
	return engine
}

func (s *HttpServer) setupRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")

	m_auth := s.container.Middlewares["auth"]
	authController := s.container.Controller.Auth
	userController := s.container.Controller.User

	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authController.Signup)
			auth.POST("/login", authController.Login)
			auth.GET("/refresh", authController.Refresh)
			auth.GET("/logout", m_auth, authController.Logout)
		}

		users := v1.Group("/users")
		{
			users.GET("/me", m_auth, userController.GetMe)
			users.GET("/:id", m_auth, userController.FindUserById)
			users.GET("/", m_auth, userController.FindUsers)
			users.PATCH("/", m_auth, userController.UpdateUser)
			users.DELETE("/", m_auth, userController.DeleteUser)
		}
	}

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusTeapot)
	})

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Route %s not found", ctx.Request.URL),
		})
	})
}
