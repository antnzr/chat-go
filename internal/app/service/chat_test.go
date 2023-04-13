package service

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

var user1 *domain.User
var user2 *domain.User

const MESSAGES_COUNT = 7

func TestCreateMessage(t *testing.T) {
	s := setup(t)
	defer s(t)

	t.Run("create message with new chats", func(t *testing.T) {
		text := gofakeit.Phrase()
		res, err := chatSrvc.CreateMessage(context.TODO(), &dto.SendMessageRequest{
			SourceUserId: user1.Id,
			TargetUserId: user2.Id,
			Text:         text,
		})

		assert.NoError(t, err)
		assert.Equal(t, text, res.Text)
	})

	t.Run("create message with existed chat", func(t *testing.T) {
		text := gofakeit.Phrase()
		res, err := chatSrvc.CreateMessage(context.TODO(), &dto.SendMessageRequest{
			SourceUserId: user1.Id,
			TargetUserId: user2.Id,
			Text:         text,
		})

		var chat int
		_ = testDbInstance.QueryRow(context.TODO(), "SELECT COUNT(*) FROM chats;").Scan(&chat)
		var messageCount int
		_ = testDbInstance.QueryRow(context.TODO(), "SELECT COUNT(*) FROM messages;").Scan(&messageCount)

		assert.NoError(t, err)
		assert.Equal(t, 1, chat)
		assert.Equal(t, 2, messageCount)
		assert.Equal(t, text, res.Text)
	})
}

func TestGetChatMessages(t *testing.T) {
	s := setup(t)
	defer s(t)
	m := makeMessages(t)
	defer m(t)

	ctx := context.TODO()
	chatId := getChatId(ctx)

	t.Run("get all chat messages", func(t *testing.T) {
		res, err := chatSrvc.FindChatMessages(ctx, &dto.FindMessagesRequest{
			ChatId: chatId,
			UserId: user1.Id,
			Limit:  20,
		})

		assert.NoError(t, err)
		assert.True(t, len(res.Docs) == MESSAGES_COUNT)
	})
	t.Run("get chat messages with next cursor", func(t *testing.T) {
		res, err := chatSrvc.FindChatMessages(ctx, &dto.FindMessagesRequest{
			ChatId: chatId,
			UserId: user1.Id,
			Limit:  3,
		})

		assert.NoError(t, err)
		assert.Empty(t, res.PrevCursor)
		assert.NotEmpty(t, res.NextCursor)
	})
	t.Run("get next chat messages with next cursor", func(t *testing.T) {
		first, err := chatSrvc.FindChatMessages(ctx, &dto.FindMessagesRequest{
			ChatId: chatId,
			UserId: user1.Id,
			Limit:  3,
		})

		res, err := chatSrvc.FindChatMessages(ctx, &dto.FindMessagesRequest{
			ChatId: chatId,
			UserId: user1.Id,
			Cursor: first.NextCursor,
			Limit:  3,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, res.PrevCursor)
		assert.NotEmpty(t, res.NextCursor)
		assert.NotEqualValues(t, first.Docs, res.Docs)
	})
}

func getChatId(ctx context.Context) int {
	searchDto := &dto.ChatSearchQuery{
		UserId: user1.Id,
	}
	searchDto.AbstractSearch = &dto.AbstractSearch{
		Limit: 1,
		Page:  1,
	}

	_, chats, _ := testStore.Chat.FindChats(ctx, *searchDto)
	return chats[0].Id
}

func setup(t *testing.T) func(tb testing.TB) {
	ctx := context.TODO()
	user1, _ = userSrvc.Signup(ctx, &dto.SignupRequest{
		Email:    gofakeit.Email(),
		Password: "password",
	})
	user2, _ = userSrvc.Signup(ctx, &dto.SignupRequest{
		Email:    gofakeit.Email(),
		Password: "password",
	})
	return func(t testing.TB) {}
}

func makeMessages(t *testing.T) func(tb testing.TB) {
	userIds := []int{user1.Id, user2.Id}
	rand.Seed(time.Now().UnixNano())
	ctx := context.TODO()
	for i := 1; i <= MESSAGES_COUNT; i++ {
		rand.Shuffle(len(userIds), func(i, j int) { userIds[i], userIds[j] = userIds[j], userIds[i] })
		_, _ = chatSrvc.CreateMessage(ctx, &dto.SendMessageRequest{
			SourceUserId: userIds[0],
			TargetUserId: userIds[1],
			Text:         gofakeit.Phrase(),
		})
	}
	return func(t testing.TB) {}
}
