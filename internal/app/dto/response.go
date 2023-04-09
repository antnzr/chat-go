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

type ChatResponse struct {
	Id            int              `json:"id"`
	Name          *string          `json:"name"`
	Description   *string          `json:"description"`
	LastMessageId *int             `json:"lastMessageId"`
	CreatedAt     time.Time        `json:"createdAt"`
	LastMessage   *MessageResponse `json:"lastMessage"`
}

type MessageResponse struct {
	Id        int       `json:"id"`
	OwnerId   int       `json:"ownerId"`
	Text      string    `json:"text"`
	ChatId    int       `json:"chatId"`
	CreatedAt time.Time `json:"createdAt"`
}
