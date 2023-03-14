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
	if u.Limit > utils.MaxLimitPerPage {
		return errs.LimitExceeded
	}
	if u.Limit < 0 || u.Page < 0 {
		return errs.BadRequest
	}
	return nil
}

type SearchResponse struct {
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	Total      int           `json:"total"`
	TotalPages int           `json:"totalPages"`
	Docs       []interface{} `json:"docs"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenDetails struct {
	Token     *string
	TokenUuid string
	UserId    int
	ExpiresIn *int64
}
