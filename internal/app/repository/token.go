package repository

import (
	"context"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tokenRepository struct {
	DB *pgxpool.Pool
}

func NewTokneRepository(db *pgxpool.Pool) domain.TokenRepository {
	return &tokenRepository{
		DB: db,
	}
}

func (tr *tokenRepository) Save(data *dto.CreateRefreshToken) (*domain.Token, error) {
	conn, err := tr.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var token domain.Token
	sqlQuery := `INSERT INTO refresh_tokens ("id", "token", "user_id")
		VALUES ($1, $2, $3)
		RETURNING "id", "token", "user_id", "created_at";`
	row := conn.QueryRow(context.Background(), sqlQuery, &data.TokenId, &data.RefreshToken, &data.UserId)

	if err := row.Scan(&token.Id, &token.RefreshToken, &token.UserId, &token.CreatedAt); err != nil {
		return &domain.Token{}, err
	}

	return &token, nil
}

func (tr *tokenRepository) DeleteByUserId(userId int) error {
	return nil
}
