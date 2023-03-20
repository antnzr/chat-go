package dto

import "time"

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type SearchResponse struct {
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	Total      int           `json:"total"`
	TotalPages int           `json:"totalPages"`
	Docs       []interface{} `json:"docs"`
}

type UserResponse struct {
	Id        int       `json:"id"`
	Email     string    `json:"email,omitempty"`
	FirstName *string   `json:"firstName,omitempty"`
	LastName  *string   `json:"lastName,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
