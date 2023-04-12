package dto

import (
	"time"
)

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type DocsResponse struct {
	Limit int           `json:"limit"`
	Docs  []interface{} `json:"docs"`
}

type PageResponse struct {
	DocsResponse
	Page       int `json:"page"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type CursorResponse struct {
	DocsResponse
	PrevCursor string `json:"prevCursor"`
	NextCursor string `json:"nextCursor"`
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
