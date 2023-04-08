package service

import (
	"os"
	"testing"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/db"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testDbInstance *pgxpool.Pool
var userSrvc domain.UserService
var messageSrvc domain.MessageService
var conf config.Config

func TestMain(m *testing.M) {
	var err error
	conf, err = config.LoadConfig("../../..")
	if err != nil || conf.Env != "test" {
		panic(err)
	}

	testDB := db.SetupTestDatabase(&conf)
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()

	buildDeps()

	os.Exit(m.Run())
}

func buildDeps() {
	userRepository := repository.NewUserRepository(testDbInstance)
	tokenRepository := repository.NewTokneRepository(testDbInstance)
	messageRepository := repository.NewMessageRepository(testDbInstance)
	store := repository.NewStore(userRepository, tokenRepository, messageRepository)

	tokenService := NewTokenService(store, conf)
	userSrvc = NewUserService(store, conf, tokenService)
	messageSrvc = NewMessageService(store, conf)
}
