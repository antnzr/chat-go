//go:build exclude
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/antnzr/chat-go/internal/app/ws"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "home.domain", "http service address")
var token = flag.String("token", "", "access jwt token")
var receiver = flag.Int("to", 0, "receiver")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/ws"}
	log.Printf("\nconnecting to: %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": []string{"Bearer " + *token}})
	if err != nil {
		log.Fatal("dial err: ", err)
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("\nread err: %v\n", err)
				}
				return
			}
			log.Printf("\nmsg: %s\n", message)
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	pong := make(chan struct{})
	conn.SetPongHandler(func(appData string) error {
		pong <- struct{}{}
		return nil
	})
	defer close(pong)

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			data, _ := json.Marshal(map[string]interface{}{
				"message": fmt.Sprintf("Hello from [%d]: %s", *receiver, t.String()),
				"receiver": receiver,
			})

			request := ws.WsEvent{
				Type:    ws.SEND_MESSAGE_EVENT,
				Payload: json.RawMessage(data),
			}

			msg, err := json.Marshal(request)
			if err != nil {
				log.Println("fail unmarshal:", err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write err:", err)
				return
			}
		case <-interrupt:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "buy buy"))
			if err != nil {
				log.Println("write close err:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(2 * time.Second):
			}
			return
		}
	}
}
