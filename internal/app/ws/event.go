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
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Receiver Client    `json:"receiver"`
	SentAt   time.Time `json:"sentAt"`
}
