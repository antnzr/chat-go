package service

import (
	"context"
	"os"
	"testing"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/db"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var testDbInstance *pgxpool.Pool
var userSrvc domain.UserService

func TestMain(m *testing.M) {
	conf, err := config.LoadConfig("../../..")
	if err != nil || conf.Env != "test" {
		panic(err)
	}

	testDB := db.SetupTestDatabase(&conf)
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()

	buildDeps()

	os.Exit(m.Run())
}

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

func buildDeps() {
	userRepository := repository.NewUserRepository(testDbInstance)
	tokenRepository := repository.NewTokneRepository(testDbInstance)
	store := repository.NewStore(userRepository, tokenRepository)

	tokenService := NewTokenService(store)
	userSrvc = NewUserService(store, tokenService)
}
