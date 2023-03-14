package repository

import (
	"context"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/jackc/pgx/v5"
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
	sqlQuery := `SELECT "id", "token", "user_id", "created_at" FROM "refresh_tokens" WHERE "id" = $1;`
	row := tr.DB.QueryRow(ctx, sqlQuery, tokenId)

	var token domain.Token
	if err := row.Scan(&token.Id, &token.RefreshToken, &token.UserId, &token.CreatedAt); err != nil {
		return nil, errs.ClarifyError(err)
	}

	return &token, nil
}

func (tr *tokenRepository) Save(ctx context.Context, data *dto.TokenDetails) (*domain.Token, error) {
	sqlQuery := `
		INSERT INTO "refresh_tokens" ("id", "token", "user_id")
		VALUES (@id, @token, @userId)
		RETURNING "id", "token", "user_id", "created_at";
	`
	args := pgx.NamedArgs{
		"id":     &data.TokenUuid,
		"token":  &data.Token,
		"userId": &data.UserId,
	}
	row := tr.DB.QueryRow(ctx, sqlQuery, args)

	var token domain.Token
	if err := row.Scan(&token.Id, &token.RefreshToken, &token.UserId, &token.CreatedAt); err != nil {
		return nil, errs.ClarifyError(err)
	}

	return &token, nil
}

func (tr *tokenRepository) DeleteByUserId(ctx context.Context, userId int) error {
	_, err := tr.DB.Exec(ctx, `DELETE FROM "refresh_tokens" WHERE "user_id" = $1;`, userId)
	if err != nil {
		return errs.ClarifyError(err)
	}
	return nil
}

func (tr *tokenRepository) DeleteToken(ctx context.Context, id string) error {
	_, err := tr.DB.Exec(ctx, `DELETE FROM "refresh_tokens" WHERE "id" = $1;`, id)
	if err != nil {
		return errs.ClarifyError(err)
	}
	return nil
}
