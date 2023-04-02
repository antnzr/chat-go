package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/antnzr/chat-go/internal/app/service"
	"github.com/gin-gonic/gin"

	_ "github.com/antnzr/chat-go/docs"
)

type userController struct {
	service service.Service
	config  config.Config
}

type UserController interface {
	GetMe(c *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	FindUsers(ctx *gin.Context)
	FindUserById(ctx *gin.Context)
}

func NewUserController(service service.Service, config config.Config) UserController {
	return &userController{service, config}
}

// GetMe godoc
// @Summary Get me
// @Description Get my user information
// @Tags User
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Success 200 {object} dto.UserResponse
// @Failure 401
// @Failure 403
// @Router /users/me [get]
func (uc *userController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	ctx.JSON(http.StatusOK, mapToUserResponse(currentUser))
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user's information
// @Tags User
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Param dto body dto.UserUpdateRequest true "User's data to update"
// @Success 200 {object} dto.UserResponse
// @Failure 400
// @Failure 401
// @Failure 403
// @Router /users [patch]
func (uc *userController) UpdateUser(ctx *gin.Context) {
	var dto dto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		_ = ctx.Error(err)
		return
	}
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	updated, err := uc.service.User.Update(context.TODO(), currentUser.Id, &dto)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapToUserResponse(updated))
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user
// @Tags User
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Router /users [delete]
func (uc *userController) DeleteUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)

	err := uc.service.User.Delete(context.TODO(), currentUser.Id)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

// FindUsers godoc
// @Summary Find users
// @Description Find users
// @Tags User
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Param limit query int  		false  "Limit per page"
// @Param page  query int  		false  "Page number"
// @Param email query string  false  "Search by email"
// @Success 200 {object} dto.SearchResponse
// @Failure 400
// @Failure 401
// @Failure 403
// @Router /users [get]
func (uc *userController) FindUsers(ctx *gin.Context) {
	var searchQuery dto.UserSearchQuery
	if err := ctx.ShouldBindQuery(&searchQuery); err != nil {
		_ = ctx.Error(err)
		return
	}
	if err := searchQuery.Validate(); err != nil {
		_ = ctx.Error(err)
		return
	}

	result, err := uc.service.User.FindAll(context.TODO(), searchQuery)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// FindUserById godoc
// @FindUserById Find user by id
// @Description Find user by id
// @Tags User
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @Security JWT
// @Param id path int true "User's id"
// @Success 200 {object} dto.UserResponse
// @Failure 400
// @Failure 401
// @Failure 403
// @Router /users/{id} [get]
func (uc *userController) FindUserById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(errs.BadRequest)
		return
	}

	user, err := uc.service.User.FindById(context.TODO(), userId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapToUserResponse(user))
}

func mapToUserResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}
