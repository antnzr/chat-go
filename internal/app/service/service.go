package service

import "github.com/antnzr/chat-go/internal/app/domain"

type Service struct {
	User  domain.UserService
	Token domain.TokenService
}

func NewService(user domain.UserService, token domain.TokenService) *Service {
	return &Service{
		User:  user,
		Token: token,
	}
}
