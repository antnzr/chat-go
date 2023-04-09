package domain

import (
	"context"
	"time"

	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/jackc/pgx/v5"
)

type Message struct {
	Id        int       `json:"id"`
	OwnerId   int       `json:"ownerId"`
	Text      string    `json:"text"`
	DialogId  int       `json:"dialogId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *Message) ScanRow(row pgx.Row) error {
	return row.Scan(
		&m.Id,
		&m.OwnerId,
		&m.DialogId,
		&m.Text,
		&m.CreatedAt,
	)
}

type MessageService interface {
	CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*Message, error)
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*Message, error)
}
