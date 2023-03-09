package controller

import (
	"net/http"

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
	var dto dto.CreateUserRequest
	if err := ctx.Bind(&dto); err != nil {
		ctx.Error(err)
		return
	}

	response, err := controller.userService.Signup(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tokens": response})
}

func (ac *authController) Login(c *gin.Context) {

}

// fmt.Errorf("%q: %w", name, ErrUserNotFound)
