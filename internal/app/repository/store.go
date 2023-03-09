package repository

import "github.com/antnzr/chat-go/internal/app/domain"

type Store struct {
	User  domain.UserRepository
	Token domain.TokenRepository
}

func NewStore(userRepo domain.UserRepository, tokenRepo domain.TokenRepository) *Store {
	return &Store{
		User:  userRepo,
		Token: tokenRepo,
	}
}
