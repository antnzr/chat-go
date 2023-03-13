package repository

import (
	"context"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
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

func (tr *tokenRepository) FindById(ctx context.Context, tokenId string) (*domain.Token, error) {
	var token domain.Token
	sqlQuery := `SELECT "id", "token", "user_id", "created_at" FROM "refresh_tokens" WHERE "id" = $1;`
	row := tr.DB.QueryRow(ctx, sqlQuery, tokenId)

	if err := row.Scan(&token.Id, &token.RefreshToken, &token.UserId, &token.CreatedAt); err != nil {
		return nil, err
	}

	return &token, nil
}

func (tr *tokenRepository) Save(ctx context.Context, data *dto.TokenDetails) (*domain.Token, error) {
	var token domain.Token
	sqlQuery := `
		INSERT INTO "refresh_tokens" ("id", "token", "user_id")
		VALUES ($1, $2, $3)
		RETURNING "id", "token", "user_id", "created_at";
	`
	row := tr.DB.QueryRow(ctx, sqlQuery, &data.TokenUuid, &data.Token, &data.UserId)

	if err := row.Scan(&token.Id, &token.RefreshToken, &token.UserId, &token.CreatedAt); err != nil {
		return nil, err
	}

	return &token, nil
}

func (tr *tokenRepository) DeleteByUserId(ctx context.Context, userId int) error {
	deleteQuery := `DELETE FROM "refresh_tokens" WHERE "user_id" = $1;`
	commandTag, err := tr.DB.Exec(ctx, deleteQuery, userId)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errs.ResourceNotFound
	}

	return nil
}

func (tr *tokenRepository) DeleteToken(ctx context.Context, id string) error {
	deleteQuery := `DELETE FROM "refresh_tokens" WHERE "id" = $1;`
	commandTag, err := tr.DB.Exec(ctx, deleteQuery, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errs.ResourceNotFound
	}

	return nil
}
