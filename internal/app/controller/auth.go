package controller

import (
	"github.com/antnzr/chat-go/internal/app/domain"
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

func (ac *authController) Signup(c *gin.Context) {

}

func (ac *authController) Login(c *gin.Context) {

}
