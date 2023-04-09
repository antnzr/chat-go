package domain

import (
	"context"

	"github.com/antnzr/chat-go/internal/app/dto"
)

type ChatService interface {
	CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*Message, error)
}

type ChatRepository interface {
	CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*Message, error)
}
