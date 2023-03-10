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

	db := db.DBPool(config)

	userRepository := repository.NewUserRepository(db)
	tokenRepository := repository.NewTokneRepository(db)
	store := repository.NewStore(userRepository, tokenRepository)

	tokenService := service.NewTokenService(store)
	userService := service.NewUserService(store, tokenService)

	authController := controller.NewAuthController(userService)
	userController := controller.NewUserController(userService)
	controller := controller.NewController(authController, userController)

	router := router.NewAppRouter(engine, *controller)
	router.Setup()

	engine.Run(fmt.Sprintf("localhost:%s", config.Port))
}
