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
var chatSrvc domain.ChatService
var conf config.Config
var testStore *repository.Store

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
	chatRepository := repository.NewChatRepository(testDbInstance)
	testStore = repository.NewStore(userRepository, tokenRepository, chatRepository)

	tokenService := NewTokenService(testStore, conf)
	userSrvc = NewUserService(testStore, conf, tokenService)
	chatSrvc = NewChatService(testStore, conf)
}
