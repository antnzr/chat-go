package domain

import (
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
	CreateTokenPair(user *User) (*dto.Tokens, error)
	DeleteByUser(userId int) error
	ValidateToken(tokenStr string, secret string) (int, error)
}

type TokenRepository interface {
	Save(data *dto.CreateRefreshToken) (*Token, error)
	DeleteByUserId(userId int) error
}
