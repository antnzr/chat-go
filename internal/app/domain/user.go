package domain

import (
	"context"
	"time"

	"github.com/antnzr/chat-go/internal/app/dto"
)

type User struct {
	Id        int        `json:"id"`
	Email     string     `json:"email,omitempty"`
	Password  string     `json:"-"`
	FirstName *string    `json:"firstName,omitempty"`
	LastName  *string    `json:"lastName,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type UserService interface {
	Signup(ctx context.Context, dto *dto.SignupRequest) (*User, error)
	Login(ctx context.Context, dto *dto.LoginRequest) (*Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
	Update(ctx context.Context, userId int, dto *dto.UserUpdateRequest) (*User, error)
	Delete(ctx context.Context, userId int) error
	FindById(ctx context.Context, id int) (*User, error)
	FindAll(ctx context.Context, searchQuery dto.UserSearchQuery) (*dto.SearchResponse, error)
}

type UserRepository interface {
	Save(ctx context.Context, dto *dto.SignupRequest) (*User, error)
	Update(ctx context.Context, userId int, dto *dto.UserUpdateRequest) (*User, error)
	FindById(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, searchQuery dto.UserSearchQuery) (int, []User, error)
}
