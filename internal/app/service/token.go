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

func (ts *tokenService) RefreshTokenPair(refreshToken string) (*dto.Tokens, error) {
	tokenDetails, err := ts.ValidateToken(refreshToken, ts.config.RefreshTokenPublicKey)
	if err != nil {
		return nil, errs.Forbidden
	}

	tokenEntity, err := ts.store.Token.FindById(tokenDetails.TokenUuid)
	if err != nil {
		return nil, errs.Forbidden
	}

	user, err := ts.store.User.FindById(tokenEntity.UserId)
	if err != nil {
		return nil, errs.Forbidden
	}

	err = ts.store.Token.DeleteToken(tokenDetails.TokenUuid)
	if err != nil {
		return nil, errs.Forbidden
	}

	return ts.CreateTokenPair(user)
}

func (ts *tokenService) CreateTokenPair(user *domain.User) (*dto.Tokens, error) {
	refreshTokenDetails, err := ts.createToken(
		user.Id,
		ts.config.RefreshTokenExpiresIn,
		ts.config.RefreshTokenPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	_, err = ts.store.Token.Save(refreshTokenDetails)
	if err != nil {
		return nil, err
	}

	accessTokenDetails, err := ts.createToken(
		user.Id,
		ts.config.AccessTokenExpiresIn,
		ts.config.AccessTokenPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	return &dto.Tokens{
		AccessToken:  *accessTokenDetails.Token,
		RefreshToken: *refreshTokenDetails.Token,
	}, nil
}

func (ts *tokenService) ValidateToken(tokenStr string, publicKey string) (*dto.TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		logger.Error(fmt.Sprintf("could not decode: %v", err))
		return nil, errs.InternalServerError
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		logger.Error(fmt.Sprintf("validate: parse key: %v", err))
		return nil, errs.InternalServerError
	}

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errs.InvalidToken
		}
		return key, nil
	})

	if err != nil {
		return nil, errs.InvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errs.InvalidToken
	}

	userId, ok := claims["sub"].(float64)
	if !ok {
		return nil, errs.InvalidToken
	}

	return &dto.TokenDetails{
		TokenUuid: fmt.Sprint(claims["jti"]),
		UserId:    int(userId),
	}, nil
}

func (ts *tokenService) DeleteByUser(userId int) error {
	if err := ts.store.Token.DeleteByUserId(userId); err != nil {
		return err
	}
	return nil
}

func (ts *tokenService) createToken(
	userId int,
	ttl time.Duration,
	privateKey string,
) (*dto.TokenDetails, error) {
	td := &dto.TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}

	now := time.Now().UTC()
	*td.ExpiresIn = now.Add(ttl).Unix()
	td.TokenUuid = uuid.New().String()
	td.UserId = userId

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		logger.Error(fmt.Sprintf("could not decode key: %v", err))
		return nil, errs.InternalServerError
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		logger.Error(fmt.Sprintf("create: parse key: %v", err))
		return nil, errs.InternalServerError
	}

	claims := make(jwt.MapClaims)
	claims["sub"] = td.UserId
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = td.ExpiresIn
	claims["jti"] = td.TokenUuid

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		logger.Error(fmt.Sprintf("create: sign token: %v", err))
		return nil, errs.InternalServerError
	}

	return td, nil
}
