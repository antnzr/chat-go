package dto

type CreateUserRequest struct {
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
}
