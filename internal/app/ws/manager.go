package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/container"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Manager struct {
	sync.RWMutex
	config    config.Config
	clients   ClientList
	container *container.Container
	handlers  map[string]WsEventHandler
}

func NewManager(
	ctx context.Context,
	container *container.Container,
	config config.Config,
) *Manager {
	manager := &Manager{
		config:    config,
		container: container,
		clients:   make(ClientList),
		handlers:  make(map[string]WsEventHandler),
	}

	manager.setupWsEventHandlers()
	return manager
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (m *Manager) ServeWs(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(*domain.User)
	if user == nil {
		logger.Debug(`unauthorized user`)
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil) //upgrade get request to websocket protocol
	if err != nil {
		logger.Error(`failed upgrade req: %v`, zap.Error(err))
		return
	}

	client := NewClient(conn, user, m)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) setupWsEventHandlers() {
	m.handlers[SEND_MESSAGE_EVENT] = SendMessage
}

func SendMessage(wsEvent WsEvent, c *Client) error {
	fmt.Println(string(wsEvent.Payload))
	var chatevent SendMessageEvent
	if err := json.Unmarshal(wsEvent.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload %v\n", err)
	}

	var broadMessage NewMessageEvent
	broadMessage.Message = chatevent.Message
	broadMessage.From = chatevent.From
	broadMessage.SentAt = time.Now()

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed unmarshall data %v\n", err)
	}

	outgoingEvent := WsEvent{
		Payload: data,
		Type:    NEW_MESSAGE_EVENT,
	}

	for client := range c.manager.clients {
		if client.chatroom == c.chatroom {
			client.egress <- outgoingEvent
		}
	}
	return nil
}

func (m *Manager) routeWsEvent(wsEvent WsEvent, c *Client) error {
	if handler, ok := m.handlers[wsEvent.Type]; ok {
		if err := handler(wsEvent, c); err != nil {
			return err
		}
		return nil
	}
	return errors.New("unknown event")
}

func (m *Manager) addClient(client *Client) {
	logger.Info(fmt.Sprintf("connect user: %s", client.user.Email))
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		logger.Info(fmt.Sprintf("disconnect user: %s", client.user.Email))
		client.connection.Close()
		delete(m.clients, client)
	}
}

// todo: check origin, origin is empty
//nolint
func (m *Manager) checkOrigin(r *http.Request) bool {
	fmt.Printf("\nOrigin: " + r.Header.Get("Origin") + "\n")
	// return slices.Contains(m.config.Origin, r.Host)
	return true
}
