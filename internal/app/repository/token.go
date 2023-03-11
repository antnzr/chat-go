package repository

import (
	"context"
	"fmt"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/antnzr/chat-go/internal/app/logger"
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

func (tr *tokenRepository) FindById(tokenId string) (*domain.Token, error) {
	conn, err := tr.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var token domain.Token
	sqlQuery := `SELECT "id", "token", "user_id", "created_at" FROM refresh_tokens WHERE id = $1`
	row := conn.QueryRow(context.Background(), sqlQuery, tokenId)

	if err := row.Scan(&token.Id, &token.RefreshToken, &token.UserId, &token.CreatedAt); err != nil {
		return nil, err
	}

	return &token, nil
}

func (tr *tokenRepository) Save(data *dto.TokenDetails) (*domain.Token, error) {
	conn, err := tr.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var token domain.Token
	sqlQuery := `INSERT INTO refresh_tokens ("id", "token", "user_id")
		VALUES ($1, $2, $3)
		RETURNING "id", "token", "user_id", "created_at";`
	row := conn.QueryRow(context.Background(), sqlQuery, &data.TokenUuid, &data.Token, &data.UserId)

	if err := row.Scan(&token.Id, &token.RefreshToken, &token.UserId, &token.CreatedAt); err != nil {
		return nil, err
	}

	return &token, nil
}

func (tr *tokenRepository) DeleteByUserId(userId int) error {
	conn, err := tr.DB.Acquire(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to acquire a database connection: %v\n", err))
		return err
	}
	defer conn.Release()

	deleteQuery := "DELETE FROM refresh_tokens WHERE user_id = $1"
	commandTag, err := conn.Exec(context.Background(), deleteQuery, userId)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errs.ResourceNotFound
	}

	return nil
}

func (tr *tokenRepository) DeleteToken(id string) error {
	conn, err := tr.DB.Acquire(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to acquire a database connection: %v\n", err))
		return err
	}
	defer conn.Release()

	deleteQuery := "DELETE FROM refresh_tokens WHERE id = $1"
	commandTag, err := conn.Exec(context.Background(), deleteQuery, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errs.ResourceNotFound
	}

	return nil
}
