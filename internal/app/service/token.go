package service

import (
	"time"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type tokenService struct {
	store  *repository.Store
	config config.Config
}

func NewTokenService(store *repository.Store) domain.TokenService {
	config, _ := config.LoadConfig(".")
	return &tokenService{store, config}
}

func (ts *tokenService) CreateTokenPair(user *domain.User) (*dto.Tokens, error) {
	tokenId := uuid.New().String()

	refreshTokenStr, err := ts.createRefreshToken(user, tokenId)
	if err != nil {
		return nil, err
	}

	tokenData := dto.CreateRefreshToken{
		TokenId:      tokenId,
		RefreshToken: refreshTokenStr,
		UserId:       user.Id,
	}

	_, err = ts.store.Token.Save(&tokenData)
	if err != nil {
		return nil, err
	}

	accessTokenStr, err := ts.createAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.Tokens{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (ts *tokenService) createRefreshToken(user *domain.User, tokenId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"jti": tokenId,
		"exp": time.Now().Add(ts.config.RefreshTokenExpiresIn).Unix(),
	})

	refreshToken, err := token.SignedString([]byte(ts.config.RefreshTokenSecret))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (ts *tokenService) createAccessToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(ts.config.AccessTokenExpiresIn).Unix(),
	})

	accessToken, err := token.SignedString([]byte(ts.config.AccessTokenSecret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
