package domain

import (
	"time"

	"github.com/jackc/pgx/v5"
)

type Message struct {
	Id        int       `json:"id"`
	OwnerId   int       `json:"ownerId"`
	Text      string    `json:"text"`
	ChatId    int       `json:"chatId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *Message) ScanRow(row pgx.Row) error {
	return row.Scan(
		&m.Id,
		&m.OwnerId,
		&m.ChatId,
		&m.Text,
		&m.CreatedAt,
	)
}
