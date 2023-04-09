package service

import (
	"context"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/antnzr/chat-go/internal/app/utils"
)

type chatService struct {
	store  *repository.Store
	config config.Config
}

func NewChatService(store *repository.Store, config config.Config) domain.ChatService {
	return &chatService{store, config}
}

func (cs *chatService) CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*domain.Message, error) {
	return cs.store.Chat.CreateMessage(ctx, dto)
}

func (cs *chatService) FindMyChats(ctx context.Context, searchQuery dto.ChatSearchQuery) (*dto.SearchResponse, error) {
	response := dto.SearchResponse{
		Page:  searchQuery.Page,
		Limit: searchQuery.Limit,
	}

	total, users, err := cs.store.Chat.FindChats(ctx, searchQuery)
	if err != nil {
		return nil, err
	}

	response.Total = total
	response.Docs = utils.ToSliceOfAny(users)
	response.TotalPages = utils.PageCount(total, searchQuery.Limit)

	return &response, nil
}
