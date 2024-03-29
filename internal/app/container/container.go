package container

import (
	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/controller"
	"github.com/antnzr/chat-go/internal/app/middleware"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Controller  controller.Controller
	Store       repository.Store
	Service     service.Service
	Middlewares map[string]gin.HandlerFunc
}

func NewContainer(config config.Config, db *pgxpool.Pool) *Container {
	userRepository := repository.NewUserRepository(db)
	tokenRepository := repository.NewTokneRepository(db)
	chatRepository := repository.NewChatRepository(db, config)
	store := repository.NewStore(userRepository, tokenRepository, chatRepository)

	tokenService := service.NewTokenService(store, config)
	userService := service.NewUserService(store, config, tokenService)
	chatService := service.NewChatService(store, config)

	serv := service.NewService(userService, tokenService, chatService)

	authController := controller.NewAuthController(*serv, config)
	userController := controller.NewUserController(*serv, config)
	chatController := controller.NewChatController(*serv, config)

	controller := controller.NewController(authController, userController, chatController)
	auth := middleware.Auth(tokenService, userService, config)

	middlewares := make(map[string]gin.HandlerFunc)
	middlewares["auth"] = auth

	return &Container{
		Store:       *store,
		Service:     *serv,
		Controller:  *controller,
		Middlewares: middlewares,
	}
}
