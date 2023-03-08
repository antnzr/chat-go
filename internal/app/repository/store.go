package repository

import "github.com/antnzr/chat-go/internal/app/domain"

type Store struct {
	User domain.UserRepository
}

func NewStore(userRepo domain.UserRepository) *Store {
	return &Store{
		User: userRepo,
	}
}
