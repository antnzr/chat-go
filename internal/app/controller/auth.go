package controller

import (
	"context"
	"net/http"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/gin-gonic/gin"
)

type authController struct {
	userService  domain.UserService
	tokenService domain.TokenService
	config       config.Config
}

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Refresh(c *gin.Context)
}

const (
	localhost    = "localhost"
	accessToken  = "accessToken"
	refreshToken = "refreshToken"
	isLoggedIn   = "isLoggedIn"
	empty        = ""
	path         = "/"
	deleteCookie = -1
	seconds      = 60
)

func NewAuthController(userService domain.UserService, tokenService domain.TokenService) AuthController {
	config, _ := config.LoadConfig(".")
	return &authController{userService, tokenService, config}
}

func (controller *authController) Signup(ctx *gin.Context) {
	var dto dto.SignupRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	err := controller.userService.Signup(context.Background(), &dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (ac *authController) Login(ctx *gin.Context) {
	var dto dto.LoginRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	tokens, err := ac.userService.Login(context.Background(), &dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	tokensResponse(ctx, tokens, &ac.config)
}

func (ac *authController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.Error(err)
	}

	if refreshToken == "" {
		ctx.Error(errs.Forbidden)
	}

	err = ac.userService.Logout(context.Background(), refreshToken)
	if err != nil {
		ctx.Error(err)
	}

	isSecure := ac.config.GinMode != gin.DebugMode
	ctx.SetCookie(accessToken, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.SetCookie(refreshToken, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.SetCookie(isLoggedIn, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.Status(http.StatusOK)
}

func (ac *authController) Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.Error(err)
	}

	if refreshToken == "" {
		ctx.Error(errs.Forbidden)
	}

	tokens, err := ac.tokenService.RefreshTokenPair(ctx, refreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	tokensResponse(ctx, tokens, &ac.config)
}

func tokensResponse(ctx *gin.Context, tokens *dto.Tokens, config *config.Config) {
	isSecure := config.GinMode != gin.DebugMode
	ctx.SetCookie(accessToken, tokens.AccessToken, config.AccessTokenMaxAge*seconds, path, localhost, isSecure, true)
	ctx.SetCookie(refreshToken, tokens.RefreshToken, config.RefreshTokenMaxAge*seconds, path, localhost, isSecure, true)
	ctx.SetCookie(isLoggedIn, "true", config.AccessTokenMaxAge*seconds, path, localhost, isSecure, false)

	ctx.JSON(http.StatusOK, gin.H{"accessToken": tokens.AccessToken})
}

// fmt.Errorf("%q: %w", name, ErrUserNotFound)
