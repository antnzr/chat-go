package service

import (
	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/antnzr/chat-go/internal/app/utils"
)

type userService struct {
	store        *repository.Store
	config       config.Config
	tokenService domain.TokenService
}

func NewUserService(store *repository.Store, tokenService domain.TokenService) domain.UserService {
	config, _ := config.LoadConfig(".")
	return &userService{store, config, tokenService}
}

func (us *userService) Signup(signupReq *dto.SignupRequest) error {
	hash := utils.HashPassword(signupReq.Password)

	signupReq.Password = hash
	_, err := us.store.User.Save(signupReq)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) Login(loginReq *dto.LoginRequest) (*dto.Tokens, error) {
	user, err := us.store.User.FindByEmail(loginReq.Email)
	if err != nil {
		return nil, errs.IncorrectCredentials
	}

	err = utils.ComparePassword(user.Password, loginReq.Password)
	if err != nil {
		return nil, errs.IncorrectCredentials
	}

	tokens, err := us.tokenService.CreateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (us *userService) Logout(refreshToken string) error {
	tokenClaims, err := us.tokenService.ValidateToken(refreshToken, us.config.RefreshTokenPublicKey)
	if err != nil {
		return errs.Forbidden
	}

	err = us.store.Token.DeleteToken(tokenClaims.TokenUuid)
	if err != nil {
		return errs.Forbidden
	}

	return nil
}

func (us *userService) GetMe(id int) (*domain.User, error) {
	user, err := us.store.User.FindById(id)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (us *userService) FindAll() ([]domain.User, error) {
	return us.store.User.FindAll()
}
