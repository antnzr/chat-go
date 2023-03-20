package domain

import (
	"context"
	"time"
)

type Token struct {
	Id           string    `json:"id"`
	UserId       int       `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenDetails struct {
	Token     *string
	TokenUuid string
	UserId    int
	ExpiresIn *int64
}

type TokenService interface {
	CreateTokenPair(ctx context.Context, user *User) (*Tokens, error)
	DeleteByUser(ctx context.Context, userId int) error
	ValidateToken(ctx context.Context, tokenStr string, secret string) (*TokenDetails, error)
	RefreshTokenPair(ctx context.Context, refreshToken string) (*Tokens, error)
}

type TokenRepository interface {
	Save(ctx context.Context, data *TokenDetails) (*Token, error)
	DeleteByUserId(ctx context.Context, userId int) error
	DeleteToken(ctx context.Context, id string) error
	FindById(ctx context.Context, tokenId string) (*Token, error)
}
