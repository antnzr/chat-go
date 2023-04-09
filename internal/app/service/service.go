package service

import "github.com/antnzr/chat-go/internal/app/domain"

type Service struct {
	User  domain.UserService
	Token domain.TokenService
	Chat  domain.ChatService
}

func NewService(
	user domain.UserService,
	token domain.TokenService,
	chat domain.ChatService,
) *Service {
	return &Service{
		User:  user,
		Token: token,
		Chat:  chat,
	}
}
