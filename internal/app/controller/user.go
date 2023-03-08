package controller

import (
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/gin-gonic/gin"
)

type userController struct {
	userService domain.UserService
}

type UserController interface {
	GetMe(c *gin.Context)
}

func NewUserController(us domain.UserService) UserController {
	return &userController{userService: us}
}

func (uc *userController) GetMe(c *gin.Context) {

}
