package app

import (
	"fmt"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/controller"
	"github.com/antnzr/chat-go/internal/app/db"
	"github.com/antnzr/chat-go/internal/app/middleware"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/antnzr/chat-go/internal/app/router"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	config, _ := config.LoadConfig(".")
	gin.SetMode(config.GinMode)

	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Origin}
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))

	db := db.DBPool(config)

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

	engine.Run(fmt.Sprintf("localhost:%s", config.Port))
}
