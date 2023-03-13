package controller

import (
	"context"
	"net/http"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/gin-gonic/gin"
)

type userController struct {
	userService domain.UserService
}

type UserController interface {
	GetMe(c *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

func NewUserController(us domain.UserService) UserController {
	return &userController{userService: us}
}

func (uc *userController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	ctx.JSON(http.StatusOK, gin.H{"user": currentUser})
}

func (uc *userController) UpdateUser(ctx *gin.Context) {
	var dto dto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	updated, err := uc.userService.Update(context.Background(), currentUser.Id, &dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": updated})
}

func (uc *userController) DeleteUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)

	err := uc.userService.Delete(context.Background(), currentUser.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}
