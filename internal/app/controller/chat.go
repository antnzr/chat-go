package controller

import (
	"net/http"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/gin-gonic/gin"
)

type chatController struct {
	service service.Service
	config  config.Config
}

type ChatController interface {
	GetMyChats(c *gin.Context)
}

func NewChatController(service service.Service, config config.Config) ChatController {
	return &chatController{service, config}
}

// GetMyChats godoc
// @Summary Get my chats
// @Description Get my chats
// @Tags Chat
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Success 200 {object} dto. --fix
// @Failure 401
// @Failure 403
// @Router /users/me [get]
func (mc *chatController) GetMyChats(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	ctx.JSON(http.StatusOK, mapToUserResponse(currentUser))
}
