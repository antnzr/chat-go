package service

import (
	"context"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/repository"
)

type messageService struct {
	store  *repository.Store
	config config.Config
}

func NewMessageService(store *repository.Store, config config.Config) domain.MessageService {
	return &messageService{store, config}
}

func (ms *messageService) CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*domain.Message, error) {
	return ms.store.Message.CreateMessage(ctx, dto)
}
