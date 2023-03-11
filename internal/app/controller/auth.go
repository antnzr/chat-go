package controller

import (
	"net/http"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/gin-gonic/gin"
)

type authController struct {
	userService domain.UserService
	config      config.Config
}

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
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

func NewAuthController(userService domain.UserService) AuthController {
	config, _ := config.LoadConfig(".")
	return &authController{userService, config}
}

func (controller *authController) Signup(ctx *gin.Context) {
	var dto dto.SignupRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	err := controller.userService.Signup(&dto)
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

	tokens, err := ac.userService.Login(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	isSecure := ac.config.GinMode != gin.DebugMode
	ctx.SetCookie(accessToken, tokens.AccessToken, ac.config.AccessTokenMaxAge*seconds, path, localhost, isSecure, true)
	ctx.SetCookie(refreshToken, tokens.RefreshToken, ac.config.RefreshTokenMaxAge*seconds, path, localhost, isSecure, true)
	ctx.SetCookie(isLoggedIn, "true", ac.config.AccessTokenMaxAge*seconds, path, localhost, isSecure, false)

	ctx.JSON(http.StatusOK, gin.H{"accessToken": tokens.AccessToken})
}

func (ac *authController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.Error(err)
	}

	if refreshToken == "" {
		ctx.Error(errs.Forbidden)
	}

	err = ac.userService.Logout(refreshToken)
	if err != nil {
		ctx.Error(err)
	}

	isSecure := ac.config.GinMode != gin.DebugMode
	ctx.SetCookie(accessToken, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.SetCookie(refreshToken, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.SetCookie(isLoggedIn, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.Status(http.StatusOK)
}

// fmt.Errorf("%q: %w", name, ErrUserNotFound)
