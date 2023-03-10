package controller

import (
	"net/http"

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

func (uc *userController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	ctx.JSON(http.StatusOK, gin.H{"user": currentUser})
}
