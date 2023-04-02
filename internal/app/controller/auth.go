package controller

import (
	"context"
	"net/http"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/gin-gonic/gin"

	_ "github.com/antnzr/chat-go/docs"
)

type authController struct {
	service service.Service
	config  config.Config
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

func NewAuthController(service service.Service, config config.Config) AuthController {
	return &authController{service, config}
}

// Signup godoc
// @Summary Signup user
// @Description Signup user
// @Param dto body dto.SignupRequest true "Signup request"
// @Tags Authentication
// @Success 200
// @Failure 400
// @Router /auth/signup [post]
func (controller *authController) Signup(ctx *gin.Context) {
	var dto dto.SignupRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		_ = ctx.Error(err)
		return
	}

	_, err := controller.service.User.Signup(context.TODO(), &dto)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusCreated)
}

// Login godoc
// @Summary Login user
// @Description Login user
// @Param dto body dto.LoginRequest true "Login request"
// @Tags Authentication
// @Header 200 {string} string accessToken
// @Header 200 {string} string refreshToken
// @Success 200 {object} dto.LoginResponse
// @Failure 401
// @Router /auth/login [post]
func (ac *authController) Login(ctx *gin.Context) {
	var dto dto.LoginRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		_ = ctx.Error(err)
		return
	}

	tokens, err := ac.service.User.Login(context.TODO(), &dto)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	tokensResponse(ctx, tokens, &ac.config)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user
// @Tags Authentication
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Success 200
// @Failure 403
// @Router /auth/logout [get]
func (ac *authController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if refreshToken == "" {
		_ = ctx.Error(errs.Forbidden)
		return
	}

	err = ac.service.User.Logout(context.TODO(), refreshToken)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	isSecure := ac.config.GinMode != gin.DebugMode
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(accessToken, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.SetCookie(refreshToken, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.SetCookie(isLoggedIn, empty, deleteCookie, path, localhost, isSecure, true)
	ctx.Status(http.StatusOK)
}

// Refresh godoc
// @Summary Refresh tokens
// @Description Refresh tokens
// @Tags Authentication
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Header 200 {string} string accessToken
// @Header 200 {string} string refreshToken
// @Success 200 {object} dto.LoginResponse
// @Failure 401
// @Failure 403
// @Router /auth/refresh [post]
func (ac *authController) Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		_ = ctx.Error(errs.Forbidden)
		return
	}

	if refreshToken == "" {
		_ = ctx.Error(errs.Forbidden)
		return
	}

	tokens, err := ac.service.Token.RefreshTokenPair(context.TODO(), refreshToken)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	tokensResponse(ctx, tokens, &ac.config)
}

func tokensResponse(ctx *gin.Context, tokens *domain.Tokens, config *config.Config) {
	isSecure := config.GinMode != gin.DebugMode
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(accessToken, tokens.AccessToken, config.AccessTokenMaxAge*seconds, path, empty, isSecure, true)
	ctx.SetCookie(refreshToken, tokens.RefreshToken, config.RefreshTokenMaxAge*seconds, path, empty, isSecure, true)
	ctx.SetCookie(isLoggedIn, "true", config.AccessTokenMaxAge*seconds, path, empty, isSecure, false)

	ctx.JSON(http.StatusOK, dto.LoginResponse{AccessToken: tokens.AccessToken})
}

// fmt.Errorf("%q: %w", name, ErrUserNotFound)
