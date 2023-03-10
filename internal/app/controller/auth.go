package controller

import (
	"net/http"
	"os"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/gin-gonic/gin"
)

type authController struct {
	userService domain.UserService
}

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
}

func NewAuthController(us domain.UserService) AuthController {
	return &authController{userService: us}
}

func (controller *authController) Signup(ctx *gin.Context) {
	var dto dto.SignupRequest
	if err := ctx.Bind(&dto); err != nil {
		ctx.Error(err)
		return
	}

	tokens, err := controller.userService.Signup(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	buildResponse(ctx, tokens)
}

func (ac *authController) Login(ctx *gin.Context) {
	var dto dto.LoginRequest
	if err := ctx.Bind(&dto); err != nil {
		ctx.Error(err)
		return
	}

	tokens, err := ac.userService.Login(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	buildResponse(ctx, tokens)
}

func buildResponse(ctx *gin.Context, tokens *dto.Tokens) {
	isSecure := os.Getenv("GIN_MODE") != "debug"
	ctx.SetCookie("jwt", tokens.RefreshToken, 60*60*24, "/", "localhost", isSecure, true)
	ctx.JSON(http.StatusOK, gin.H{"accessToken": tokens.AccessToken})
}

// fmt.Errorf("%q: %w", name, ErrUserNotFound)
