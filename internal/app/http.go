package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/container"
	"github.com/antnzr/chat-go/internal/app/middleware"
	"github.com/antnzr/chat-go/internal/app/ws"
	"github.com/antnzr/chat-go/internal/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	ginzap "github.com/gin-contrib/zap"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpServer struct {
	config    config.Config
	container *container.Container
	http      *http.Server
}

func NewHttpServer(
	config config.Config,
	container *container.Container,
) *HttpServer {
	return &HttpServer{
		config:    config,
		container: container,
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	engine := s.setup(ctx)

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

func (s *HttpServer) setup(ctx context.Context) *gin.Engine {
	gin.SetMode(s.config.GinMode)

	engine := gin.New()
	_ = engine.SetTrustedProxies(nil)
	engine.Use(middleware.ErrorHandler())

	engine.Use(ginzap.Ginzap(logger.GetLogger(), time.RFC3339Nano, true))
	engine.Use(ginzap.RecoveryWithZap(logger.GetLogger(), true))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = s.config.Origin
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))
	s.setupRoutes(ctx, engine)

	return engine
}

func (s *HttpServer) setupRoutes(ctx context.Context, engine *gin.Engine) {
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

	manager := ws.NewManager(ctx, s.container, s.config)
	engine.GET("/ws", m_auth, manager.ServeWs)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/health", s.health)
	engine.NoRoute(s.noRoute)
}

func (s *HttpServer) noRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": fmt.Sprintf("Route %s not found", ctx.Request.URL),
	})
}

func (s *HttpServer) health(ctx *gin.Context) {
	ctx.Status(http.StatusTeapot)
}
