package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/antnzr/chat-go/config"
	_ "github.com/antnzr/chat-go/docs"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/gin-gonic/gin"
)

type chatController struct {
	service service.Service
	config  config.Config
}

type ChatController interface {
	GetMyChats(ctx *gin.Context)
	GetChatMessages(ctx *gin.Context)
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
// @Success 200 {object} dto.PageResponse
// @Failure 401
// @Failure 403
// @Router /chats [get]
func (ch *chatController) GetMyChats(ctx *gin.Context) {
	var searchQuery dto.ChatSearchQuery
	if err := ctx.ShouldBindQuery(&searchQuery); err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := searchQuery.Validate(); err != nil {
		_ = ctx.Error(err)
		return
	}

	currentUser := ctx.MustGet("currentUser").(*domain.User)
	searchQuery.UserId = currentUser.Id

	result, err := ch.service.Chat.FindMyChats(context.TODO(), searchQuery)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetChatMessages godoc
// @Summary Get chat messages
// @Description Get chat messages
// @Tags Chat
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Param chatId  path  int    true   "Chat Id"
// @Param limit   query int  	 false  "Limit per page default 20"
// @Param cursor  query string false  "Previous or next result set"
// @Success 200 {object} dto.CursorResponse
// @Failure 401
// @Failure 403
// @Router /chats/{chatId}/messages [get]
func (ch *chatController) GetChatMessages(ctx *gin.Context) {
	chatId, err := strconv.Atoi(ctx.Param("chatId"))
	if err != nil {
		_ = ctx.Error(errs.BadRequest)
		return
	}

	searchQuery := new(dto.FindMessagesRequest)
	if err := ctx.ShouldBindQuery(&searchQuery); err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := searchQuery.Validate(); err != nil {
		_ = ctx.Error(err)
		return
	}

	currentUser := ctx.MustGet("currentUser").(*domain.User)
	searchQuery.UserId = currentUser.Id
	searchQuery.ChatId = chatId

	res, err := ch.service.Chat.FindChatMessages(context.TODO(), searchQuery)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
