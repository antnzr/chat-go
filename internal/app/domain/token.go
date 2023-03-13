package domain

import (
	"context"
	"time"

	"github.com/antnzr/chat-go/internal/app/dto"
)

type Token struct {
	Id           string    `json:"id"`
	UserId       int       `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	CreatedAt    time.Time `json:"createdAt"`
}

type TokenService interface {
	CreateTokenPair(ctx context.Context, user *User) (*dto.Tokens, error)
	DeleteByUser(ctx context.Context, userId int) error
	ValidateToken(ctx context.Context, tokenStr string, secret string) (*dto.TokenDetails, error)
	RefreshTokenPair(ctx context.Context, refreshToken string) (*dto.Tokens, error)
}

type TokenRepository interface {
	Save(ctx context.Context, data *dto.TokenDetails) (*Token, error)
	DeleteByUserId(ctx context.Context, userId int) error
	DeleteToken(ctx context.Context, id string) error
	FindById(ctx context.Context, tokenId string) (*Token, error)
}
