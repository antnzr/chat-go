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

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": []string{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA0MzI0MDksImlhdCI6MTY4MDQzMTUwOSwianRpIjoiMjkyMmY2OTEtN2ExNC00N2Q2LWFjOTMtNDA1NmYwNjk5MDU3IiwibmJmIjoxNjgwNDMxNTA5LCJzdWIiOjF9.rfa478p0FCJfFx_V2TonX4hj-BZ_okfeWPlYACCU7GJNHBn3ym6Iz0yPlZ-uyUBcwyH3hV2KkG13losFutGr_Nx_n_DcLW0RFXttx8QiXE9zA50vZW45fBpQQqGP8gBb3yIZncDg-HksbDTT9fYIfGG_xmXsZrXF-uWAYcQTTaoCophf05pwvTmIoPToj4iK6DCMzUQCsg75GyenoGJ1a3SNWlH_ezcSmrrZIk2htaEJj-zvDHoJpbypDJLPGrbJI0I4i8azjSHiDUf_4MaFmgZR_MxeNE-qXZv0pXhvFe4c9tbiMa96WMBwvLT2mkZA9g_kxoDsOaDS-HeCow8kuhFRzLyBWoawNVtFdAb-PZ3O-klV6QCX1EqtUyiG1kHKxKYN7IZatqsZhFOZhu2VmR8AMqXsly2uQIXzMBVkHaSaof07chMAxOiGK6eZAVk2p0Es3dQWiJ12WL8Z9Is-q2b22E01hmwiPi4h59WOX44KLgr_wLLlAIyEKg-85x33IEwgys24CSti-NTNCy6peTFnD5DWZXZwhhbAPEqZU-mChbAUgASbkLJyt5MHwgeJ7I674GkHnXHmA8OgUFjLMRhd8jE6WqNL3jkxXoKnNxsiw4JfZ4B5RnQrhhH3WHfuimgymZpvpnOB7KsCu2zzPmLIyvThwItP-YOz1kWTIgw"}})
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
