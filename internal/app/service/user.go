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

type userService struct {
	store        *repository.Store
	config       config.Config
	tokenService domain.TokenService
}

func NewUserService(
	store *repository.Store,
	config config.Config,
	tokenService domain.TokenService,
) domain.UserService {
	return &userService{store, config, tokenService}
}

func (us *userService) Signup(ctx context.Context, signupReq *dto.SignupRequest) (*domain.User, error) {
	hash := utils.HashPassword(signupReq.Password)

	signupReq.Password = hash
	user, err := us.store.User.Save(ctx, signupReq)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) Login(ctx context.Context, loginReq *dto.LoginRequest) (*domain.Tokens, error) {
	user, err := us.store.User.FindByEmail(ctx, loginReq.Email)
	if err != nil {
		return nil, errs.IncorrectCredentials
	}

	err = utils.ComparePassword(user.Password, loginReq.Password)
	if err != nil {
		return nil, errs.IncorrectCredentials
	}

	tokens, err := us.tokenService.CreateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (us *userService) Logout(ctx context.Context, refreshToken string) error {
	tokenClaims, err := us.tokenService.ValidateToken(ctx, refreshToken, us.config.RefreshTokenPublicKey)
	if err != nil {
		return errs.Forbidden
	}

	err = us.store.Token.DeleteToken(ctx, tokenClaims.TokenUuid)
	if err != nil {
		return errs.Forbidden
	}

	return nil
}

func (us *userService) FindById(ctx context.Context, id int) (*domain.User, error) {
	user, err := us.store.User.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) Update(ctx context.Context, userId int, dto *dto.UserUpdateRequest) (*domain.User, error) {
	user, err := us.store.User.Update(ctx, userId, dto)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) Delete(ctx context.Context, userId int) error {
	return us.store.User.Delete(ctx, userId)
}

func (us *userService) FindAll(ctx context.Context, searchQuery dto.UserSearchQuery) (*dto.PageResponse, error) {
	total, docs, err := us.store.User.FindAll(ctx, searchQuery)
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
