package repository

import (
	"github.com/antnzr/chat-go/internal/app/domain"
)

type Store struct {
	User  domain.UserRepository
	Token domain.TokenRepository
	Chat  domain.ChatRepository
}

func NewStore(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	chatRepo domain.ChatRepository,
) *Store {
	return &Store{
		User:  userRepo,
		Token: tokenRepo,
		Chat:  chatRepo,
	}
}
