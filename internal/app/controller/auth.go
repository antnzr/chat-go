package controller

import (
	"net/http"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/gin-gonic/gin"
)

type authController struct {
	userService domain.UserService
	config      config.Config
}

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
}

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

	isSecure := ac.config.GinMode != "debug"
	ctx.SetCookie("accessToken", tokens.AccessToken, ac.config.AccessTokenMaxAge*60, "/", "localhost", isSecure, true)
	ctx.SetCookie("refreshToken", tokens.RefreshToken, ac.config.RefreshTokenMaxAge*60, "/", "localhost", isSecure, true)
	ctx.SetCookie("loggedIn", "true", ac.config.AccessTokenMaxAge*60, "/", "localhost", isSecure, false)

	ctx.JSON(http.StatusOK, gin.H{"accessToken": tokens.AccessToken})
}

// fmt.Errorf("%q: %w", name, ErrUserNotFound)
