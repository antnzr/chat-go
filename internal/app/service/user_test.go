package service

import (
	"context"
	"testing"

	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestSignupUser(t *testing.T) {
	email := gofakeit.Email()
	user, err := userSrvc.Signup(context.TODO(), &dto.SignupRequest{
		Email:    email,
		Password: "password",
	})

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
}

func TestFindUser(t *testing.T) {
	email := gofakeit.Email()
	signed, _ := userSrvc.Signup(context.TODO(), &dto.SignupRequest{
		Email:    email,
		Password: "password",
	})

	user, err := userSrvc.FindById(context.TODO(), signed.Id)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, signed.Id, user.Id)
}

func TestUpdateUser(t *testing.T) {
	email := gofakeit.Email()
	firstName := gofakeit.FirstName()
	signed, _ := userSrvc.Signup(context.TODO(), &dto.SignupRequest{
		Email:     email,
		FirstName: firstName,
		Password:  "password",
	})

	lastName := gofakeit.LastName()
	updated, err := userSrvc.Update(context.TODO(), signed.Id, &dto.UserUpdateRequest{
		LastName: lastName,
	})

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, lastName, *updated.LastName)
	assert.Equal(t, firstName, *signed.FirstName)
	assert.Equal(t, signed.Id, updated.Id)
}
