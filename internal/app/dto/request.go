package dto

import (
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/antnzr/chat-go/internal/app/utils"
)

type SignupRequest struct {
	Email     string `json:"email,omitempty" binding:"required"`
	Password  string `json:"password,omitempty" binding:"required"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type UserUpdateRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type UserSearchQuery struct {
	Limit int     `form:"limit,default=20"`
	Page  int     `form:"page,default=1"`
	Email *string `form:"email"`
}

func (u *UserSearchQuery) Validate() error {
	if u.Limit > utils.MAX_LIMIT_PER_PAGE {
		return errs.LimitExceeded
	}
	if u.Limit < 0 || u.Page < 0 {
		return errs.BadRequest
	}
	return nil
}
