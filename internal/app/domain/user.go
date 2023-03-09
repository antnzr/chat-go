package domain

import (
	"time"

	"github.com/antnzr/chat-go/internal/app/dto"
)

type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	FirstName *string   `json:"firstName,omitempty"`
	LastName  *string   `json:"lastName,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserService interface {
	Signup(dto *dto.CreateUserRequest) (*dto.Tokens, error)
	GetMe(id int) (*User, error)
	FindAll() ([]User, error)
}

type UserRepository interface {
	Save(dto *dto.CreateUserRequest) (*User, error)
	FindById(id int) (*User, error)
	FindAll() ([]User, error)
	Delete(id int) error
}
