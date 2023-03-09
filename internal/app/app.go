package app

import (
	"fmt"
	"os"

	"github.com/antnzr/chat-go/internal/app/controller"
	"github.com/antnzr/chat-go/internal/app/db"
	"github.com/antnzr/chat-go/internal/app/logger"
	"github.com/antnzr/chat-go/internal/app/middleware"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/antnzr/chat-go/internal/app/router"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/antnzr/chat-go/internal/pkg"
	"github.com/gin-gonic/gin"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	gin.SetMode(os.Getenv("GIN_MODE"))
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler())

	pkg.LoadEnvVars()
	db := db.DBPool()

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

	port := os.Getenv("PORT")
	engine.Run(fmt.Sprintf("localhost:%s", port))
	logger.Info(fmt.Sprintf("App is running on port: %s", port))
}
