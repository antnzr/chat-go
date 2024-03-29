package domain

import (
	"context"

	"github.com/antnzr/chat-go/internal/app/dto"
)

type ChatService interface {
	CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*dto.MessageResponse, error)
	FindMyChats(ctx context.Context, searchQuery dto.ChatSearchQuery) (*dto.PageResponse, error)
	FindChatMessages(ctx context.Context, searchQuery *dto.FindMessagesRequest) (*dto.CursorResponse, error)
}

type ChatRepository interface {
	CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*Message, error)
	FindChats(ctx context.Context, searchQuery dto.ChatSearchQuery) (int, []dto.ChatResponse, error)
	FindChatMessages(ctx context.Context, searchQuery *dto.FindMessagesRequest) ([]dto.MessageResponse, error)
}
