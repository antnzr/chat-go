package service

import (
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/repository"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	store        *repository.Store
	tokenService domain.TokenService
}

func NewUserService(store *repository.Store, tokenService domain.TokenService) domain.UserService {
	return &userService{store: store, tokenService: tokenService}
}

func (us *userService) Signup(requestData *dto.CreateUserRequest) (*dto.Tokens, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), 10)
	if err != nil {
		return nil, err
	}

	requestData.Password = string(hash)
	user, err := us.store.User.Save(requestData)
	if err != nil {
		return nil, err
	}

	tokens, err := us.tokenService.CreateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (us *userService) GetMe(id int) (*domain.User, error) {
	return us.store.User.FindById(id)
}

func (us *userService) FindAll() ([]domain.User, error) {
	return us.store.User.FindAll()
}
