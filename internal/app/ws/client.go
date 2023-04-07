package ws

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/pkg/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[int]*Client

type Client struct {
	connection *websocket.Conn
	user       *domain.User
	manager    *Manager
	chatroom   string
	egress     chan WsEvent // is used to avoid concurrent writes on the WebSocket
}

func NewClient(conn *websocket.Conn, user *domain.User, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		user:       user,
		egress:     make(chan WsEvent),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	c.connection.SetReadLimit(5121)
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logger.Err(zap.Error(err))
	}
	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Debug("close: %v", zap.Error(err))
				return
			}
			logger.Error("read err: %v", zap.Error(err))
			return
		}

		var request WsEvent

		if err := json.Unmarshal(payload, &request); err != nil {
			logger.Error("error marshalling event: %v\n", zap.Error(err))
		}

		if err := c.manager.routeWsEvent(request, c); err != nil {
			logger.Error("error route wsEvent: %v\n", zap.Error(err))
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					logger.Error("connection is closed: ", zap.Error(err))
					return
				}
			}

			data, err := json.Marshal(message)
			if err != nil {
				logger.Err(zap.Error(err))
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				logger.Error("failed send message: %v", zap.Error(err))
			}
		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				if strings.Contains(err.Error(), "websocket: close sent") {
					return
				}
				logger.Error("ping err: %v\n", zap.Error(err))
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
