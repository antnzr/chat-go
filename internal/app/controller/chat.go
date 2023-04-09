package controller

import (
	"context"
	"net/http"

	"github.com/antnzr/chat-go/config"
	_ "github.com/antnzr/chat-go/docs"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
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
// @Param limit query int  		false  "Limit per page default 20"
// @Param page  query int  		false  "Page number default 1"
// @Success 200 {object} dto.SearchResponse
// @Failure 401
// @Failure 403
// @Router /chats [get]
func (mc *chatController) GetMyChats(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	var searchQuery dto.ChatSearchQuery
	if err := ctx.ShouldBindQuery(&searchQuery); err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := searchQuery.Validate(); err != nil {
		_ = ctx.Error(err)
		return
	}

	searchQuery.UserId = currentUser.Id
	result, err := mc.service.Chat.FindMyChats(context.TODO(), searchQuery)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
