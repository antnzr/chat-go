package repository

import (
	"github.com/antnzr/chat-go/internal/app/domain"
)

type Store struct {
	User    domain.UserRepository
	Token   domain.TokenRepository
	Message domain.MessageRepository
}

func NewStore(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	messageRepo domain.MessageRepository,
) *Store {
	return &Store{
		User:    userRepo,
		Token:   tokenRepo,
		Message: messageRepo,
	}
}
