package service

import "github.com/antnzr/chat-go/internal/app/domain"

type Service struct {
	User    domain.UserService
	Token   domain.TokenService
	Message domain.MessageService
}

func NewService(
	user domain.UserService,
	token domain.TokenService,
	message domain.MessageService,
) *Service {
	return &Service{
		User:    user,
		Token:   token,
		Message: message,
	}
}
