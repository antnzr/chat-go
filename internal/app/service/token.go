package service

import (
	"time"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
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

func (ts *tokenService) ValidateToken(tokenStr string, secret string) (int, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.InvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, errs.InvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return 0, errs.InvalidToken
	}

	userId, ok := claims["sub"].(float64)
	if !ok {
		return 0, errs.InvalidToken
	}

	return int(userId), nil
}

func (ts *tokenService) DeleteByUser(userId int) error {
	if err := ts.store.Token.DeleteByUserId(userId); err != nil {
		return err
	}
	return nil
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
