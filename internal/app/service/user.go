package service

import (
	"sync"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/repository"
)

type userService struct {
	store *repository.Store
}

var once sync.Once
var instance *userService

func NewUserService(store *repository.Store) domain.UserService {
	once.Do(func() { instance = &userService{store} })
	return instance
}

func (us *userService) Signup(dto *dto.CreateUserRequest) (*domain.User, error) {
	return us.store.User.Save(dto)
}

func (us *userService) GetMe(id int) (*domain.User, error) {
	return us.store.User.FindById(id)
}

func (us *userService) FindAll() ([]domain.User, error) {
	return us.store.User.FindAll()
}
