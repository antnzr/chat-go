package ws

import (
	"encoding/json"
	"time"
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
	SendMessageEvent
	From     int       `json:"from"`
	Receiver Client    `json:"receiver"`
	SentAt   time.Time `json:"sentAt"`
}
