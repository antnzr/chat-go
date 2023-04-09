package service

import (
	"context"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/repository"
)

type chatService struct {
	store  *repository.Store
	config config.Config
}

func NewChatService(store *repository.Store, config config.Config) domain.ChatService {
	return &chatService{store, config}
}

func (ms *chatService) CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*domain.Message, error) {
	return ms.store.Chat.CreateMessage(ctx, dto)
}
