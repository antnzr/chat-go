package ws

import (
	"encoding/json"

	"github.com/antnzr/chat-go/internal/app/domain"
)

type WsEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WsEventHandler func(wsEvent WsEvent, c *Client) error

const (
	SEND_MESSAGE_EVENT = "send_message"
	NEW_MESSAGE_EVENT  = "new_message"
	ERROR_EVENT        = "error"
)

type ErrorEvent struct {
	Message string `json:"message"`
}

type SendMessageEvent struct {
	Message  string `json:"message"`
	Receiver int    `json:"receiver"`
}

type NewMessageEvent struct {
	Message domain.Message `json:"message"`
}
