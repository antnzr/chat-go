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

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/ws"}
	log.Printf("\nconnecting to: %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": []string{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA0MjIxNTEsImlhdCI6MTY4MDQyMTI1MSwianRpIjoiYzQ2YTU5ZTktZWU2Yy00NzNlLThlNjItZGRkM2QyMzkwNTM3IiwibmJmIjoxNjgwNDIxMjUxLCJzdWIiOjF9.wNI2mw-7OYLvnE1ftqqDpV8471z0anCooYa10NukZ7RhYtyjtT2hwEgpNPCTUZEjaxd-A16tTKxWl_lvhGFxD1fPof689eenQIWW2hR-WlltYKYA0YCQPDR43xEMe_CF8wxzmeiLVwpSVZczgzLni2sMxnlHJR_SFAe-5WpzyTAXZpzaGe-XLdQKWw5_kSYQlIHNgrxowHr6gyZqCK6hOVPjQJDrfliqfuX612AsxTmQisO_DxH8O6b5RkzYOn19QBmiF1HqVm9ciDhI7YOXQIIgfJ7o78Hh5QaAVj1FpVus0-O-5kAXRIB3ur_y_YhHMEGxeBqzOXV2jGSIr7ZetDn-GeGrJx89jhKWZhKxTDUhF_EqKfvfc4wqcUg5uhgHUwTBuBxgffmDYcBXMuv25jbsBqWZlhfVgY6arm9B5V4D-swl59b9xV2q034Gk0GX5HQiynnEgiVYxl4ev7BONHQUJTFLBMQFehRAVv6qPKZ0L_3MFGqkbV7wTwDS5d8FOE3-zO5korFV-0mcgG00WQ8aPdXRwPLmIO-BxOF5iX-NVvc0myBjhzN_FQb2MVtAHt3MYdKZdm9tcddNck44b8mZoqXi-L6oGcG-DZjcHy1u1DlDQRzTFtsq4vFsK-KhDgOboPnvamcfqRg-lr_13rBR3du88JzUIJuA7EGY49o"}})
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

	ticker := time.NewTicker(time.Second * 3)
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
			data, _ := json.Marshal(map[string]string{"message": fmt.Sprintf("Hello muchacho: %s", t.String())})

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
