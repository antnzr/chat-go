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

type Search interface {
	Validate() error
}

type AbstractSearch struct {
	Limit int `form:"limit,default=20"`
	Page  int `form:"page,default=1"`
	Search
}

func (as *AbstractSearch) Validate() error {
	if as.Limit > utils.MAX_LIMIT_PER_PAGE {
		return errs.LimitExceeded
	}
	if as.Limit < 0 || as.Page < 0 {
		return errs.BadRequest
	}
	return nil
}

type UserSearchQuery struct {
	*AbstractSearch
	Email *string `form:"email"`
}

type ChatSearchQuery struct {
	*AbstractSearch
	UserId int `json:"-"`
}

type SendMessageRequest struct {
	SourceUserId int
	TargetUserId int
	Text         string
}

type FindMessagesRequest struct {
	Limit         int          `form:"limit,default=20"`
	SortOrder     string       `form:"sortOrder,default=desc"`
	Cursor        string       `form:"cursor"`
	UserId        int          `json:"-"`
	ChatId        int          `json:"-"`
	DecodedCursor utils.Cursor `json:"-"`
}

func (fm *FindMessagesRequest) Validate() error {
	if fm.Limit > utils.MAX_LIMIT_PER_PAGE {
		return errs.LimitExceeded
	}
	if fm.Limit < 0 {
		return errs.BadRequest
	}
	return nil
}
