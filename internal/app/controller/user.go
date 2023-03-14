package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/gin-gonic/gin"
)

type userController struct {
	userService domain.UserService
}

type UserController interface {
	GetMe(c *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	FindUsers(ctx *gin.Context)
	FindUserById(ctx *gin.Context)
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

func (uc *userController) FindUsers(ctx *gin.Context) {
	var searchQuery dto.UserSearchQuery
	if err := ctx.ShouldBindQuery(&searchQuery); err != nil {
		ctx.Error(err)
		return
	}
	if err := searchQuery.Validate(); err != nil {
		ctx.Error(err)
		return
	}

	result, err := uc.userService.FindAll(context.Background(), searchQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": &result})
}

func (uc *userController) FindUserById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(errs.BadRequest)
		return
	}

	user, err := uc.userService.FindById(context.Background(), userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": &user})
}
