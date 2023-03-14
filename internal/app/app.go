package app

import (
	"context"
	"fmt"
	"time"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/controller"
	"github.com/antnzr/chat-go/internal/app/db"
	"github.com/antnzr/chat-go/internal/app/logger"
	"github.com/antnzr/chat-go/internal/app/middleware"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/antnzr/chat-go/internal/app/router"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	config, _ := config.LoadConfig(".")
	gin.SetMode(config.GinMode)

	engine := gin.New()
	engine.SetTrustedProxies(nil)
	engine.Use(middleware.ErrorHandler())

	engine.Use(ginzap.Ginzap(logger.GetLogger(), time.RFC3339Nano, true))
	engine.Use(ginzap.RecoveryWithZap(logger.GetLogger(), true))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Origin
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))

	db, err := db.DBPool(context.Background(), config)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	tokenRepository := repository.NewTokneRepository(db)
	store := repository.NewStore(userRepository, tokenRepository)

	tokenService := service.NewTokenService(store)
	userService := service.NewUserService(store, tokenService)

	authController := controller.NewAuthController(userService, tokenService)
	userController := controller.NewUserController(userService)
	controller := controller.NewController(authController, userController)

	auth := middleware.Auth(tokenService, userService, config)

	router := router.NewAppRouter(engine, controller, auth)
	router.Setup()

	logger.Fatal("%v", zap.Error(engine.Run(fmt.Sprintf("localhost:%s", config.Port))))
}
