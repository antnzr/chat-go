package service

import (
	"context"
	"log"
	"testing"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

var user1 *domain.User
var user2 *domain.User

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
	return func(t testing.TB) {
		log.Println("teardown test")
	}
}

func TestCreateMessage(t *testing.T) {
	s := setup(t)
	defer s(t)

	text := gofakeit.Phrase()
	res, err := messageSrvc.SendMessage(context.TODO(), &dto.SendMessageRequest{
		SourceUserId: user1.Id,
		TargetUserId: user2.Id,
		Text:         text,
	})

	assert.NoError(t, err)
	assert.Equal(t, text, res.Text)
}
