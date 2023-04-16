package service

import (
	"context"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
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

func (cs *chatService) CreateMessage(ctx context.Context, query *dto.SendMessageRequest) (*dto.MessageResponse, error) {
	msg, err := cs.store.Chat.CreateMessage(ctx, query)
	if err != nil {
		return nil, err
	}

	msgResponse := &dto.MessageResponse{
		Id:        msg.Id,
		OwnerId:   msg.OwnerId,
		Text:      msg.Text,
		ChatId:    msg.ChatId,
		CreatedAt: msg.CreatedAt,
	}

	return msgResponse, nil
}

func (cs *chatService) FindMyChats(ctx context.Context, searchQuery dto.ChatSearchQuery) (*dto.PageResponse, error) {
	total, docs, err := cs.store.Chat.FindChats(ctx, searchQuery)
	if err != nil {
		return nil, err
	}

	res := new(dto.PageResponse)
	res.Total = total
	res.TotalPages = utils.PageCount(total, searchQuery.Limit)
	res.Page = searchQuery.Page
	res.Limit = searchQuery.Limit
	res.Docs = utils.ToSliceOfAny(docs)

	return res, nil
}

func (cs *chatService) FindChatMessages(ctx context.Context, query *dto.FindMessagesRequest) (*dto.CursorResponse, error) {
	isFirstPage := query.Cursor == ""
	isPointNext := false

	if query.Cursor != "" {
		decodedCursor, err := utils.DecodeCursor(query.Cursor)
		if err != nil {
			return nil, errs.BadRequest
		}
		isPointNext = decodedCursor.IsPointNext
		query.DecodedCursor = *decodedCursor
	}

	messages, err := cs.store.Chat.FindChatMessages(ctx, query)
	if err != nil {
		return nil, err
	}

	hasMore := len(messages) > int(query.Limit)
	if hasMore {
		messages = messages[:query.Limit]
	}
	if !isFirstPage && !isPointNext {
		messages = utils.Reverse(messages)
	}

	cursors := getCursors(isFirstPage, hasMore, query.Limit, messages, isPointNext)

	res := new(dto.CursorResponse)
	res.Limit = query.Limit
	res.Docs = utils.ToSliceOfAny(messages)
	if cursors != nil {
		res.PrevCursor = cursors.PrevCursor
		res.NextCursor = cursors.NextCursor
	}

	return res, nil
}

func getCursors(
	isFirstPage bool,
	hasMore bool,
	limit int,
	messages []dto.MessageResponse,
	isPointNext bool,
) *utils.Cursors {
	nextCur := new(utils.Cursor)
	prevCur := new(utils.Cursor)

	if isFirstPage {
		if hasMore {
			nextCur := utils.NewCursor(messages[limit-1].Id, true)
			return utils.NewCursors(nextCur, nil)
		}
		return nil
	}

	if isPointNext {
		if hasMore {
			nextCur = utils.NewCursor(messages[limit-1].Id, true)
		}
		prevCur = utils.NewCursor(messages[0].Id, false)
		return utils.NewCursors(nextCur, prevCur)
	}

	if len(messages) > 1 {
		nextCur = utils.NewCursor(messages[len(messages)-1].Id, true)
	}

	if hasMore {
		prevCur = utils.NewCursor(messages[0].Id, false)
	}

	return utils.NewCursors(nextCur, prevCur)
}
