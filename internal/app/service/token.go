package service

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/antnzr/chat-go/internal/app/logger"
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
	refreshTokenJti := uuid.New().String()

	refreshToken, err := ts.createToken(
		user.Id,
		ts.config.RefreshTokenExpiresIn,
		ts.config.RefreshTokenPrivateKey,
		refreshTokenJti,
	)
	if err != nil {
		return nil, err
	}

	tokenData := dto.CreateRefreshToken{
		TokenId:      refreshTokenJti,
		RefreshToken: refreshToken,
		UserId:       user.Id,
	}

	_, err = ts.store.Token.Save(&tokenData)
	if err != nil {
		return nil, err
	}

	accessToken, err := ts.createToken(
		user.Id,
		ts.config.AccessTokenExpiresIn,
		ts.config.AccessTokenPrivateKey,
		"",
	)
	if err != nil {
		return nil, err
	}

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ts *tokenService) ValidateToken(tokenStr string, publicKey string) (int, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return 0, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return 0, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errs.InvalidToken
		}
		return key, nil
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

func (ts *tokenService) createToken(
	payload interface{},
	ttl time.Duration,
	privateKey string,
	jti string,
) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		logger.Error(fmt.Sprintf("could not decode key: %v", err))
		return "", errs.InternalServerError
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		logger.Error(fmt.Sprintf("create: parse key: %v", err))
		return "", errs.InternalServerError
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(ttl).Unix()

	if jti != "" {
		claims["jti"] = jti
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		logger.Error(fmt.Sprintf("create: sign token: %v", err))
		return "", errs.InternalServerError
	}

	return token, nil
}
