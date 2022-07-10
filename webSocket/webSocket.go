package webSocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jarrancarr/ChexxServer/store"
)

type WsHandler struct{}

//var clients = make(map[*websocket.Conn]bool)
var socket = make(map[uint]*websocket.Conn)

func (wsh WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Upgrading HTTP Connection to websocket connection
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  256,
		WriteBufferSize: 256,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, upgradeErr := upgrader.Upgrade(w, r, nil)
	if upgradeErr != nil {
		log.Printf("error upgrading %s", upgradeErr)
		return
	}

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			msg := &struct {
				Type    string `json:"type"`
				Token   string `json:"token"`
				Epoc    int64  `json:"epoc"`
				Message string `json:"message"`
			}{}
			read := bytes.NewReader(message)
			err = json.NewDecoder(read).Decode(msg)
			if err != nil {
				log.Println("decode:", err)
				break
			}
			fmt.Printf("type: %s  token: %s  Date: %v  Message: %s\n", msg.Type, msg.Token, msg.Epoc, msg.Message)
			//conn.WriteMessage(1, []byte("I hear you"))
			if store.SessionMap[msg.Token] != nil {
				switch msg.Type {
				case "ping":
					store.Sessions()[msg.Token].Inbox <- fmt.Sprintf("bounce||%d", time.Now().Unix())
				case "login":
					store.SessionMap[msg.Token].WsConn = conn
					go wsDataQueue(msg.Token)
				}
			} else {
				log.Println("No token or no session")
			}
		}
	}()
}

func wsDataQueue(token string) {
	if store.SessionMap[token].WsConn != nil {
		for {
			d := <-store.SessionMap[token].Inbox
			switch d.(type) {
			case string:
				pair := strings.Split(d.(string), "|||")
				packet := "{"
				for i := 0; i < len(pair); i += 1 {
					if i > 0 {
						packet += ", "
					}
					spl := strings.Split(pair[i], "||")
					packet += "\"" + spl[0] + "\":\"" + spl[1] + "\""
					//packet += spl[0] + ":\"" + spl[1] + "\""
				}
				packet += "}"
				store.SessionMap[token].WsConn.WriteMessage(1, []byte(packet))
			// case *store.Message:
			// 	msg, _ := json.Marshal(d)
			// 	store.SessionMap[token].WsConn.WriteMessage(1, []byte("{\"chat\":"+string(msg)+"}"))
			case *store.Match:
				match, _ := json.Marshal(d)
				store.SessionMap[token].WsConn.WriteMessage(1, []byte("{\"match\":"+string(match)+"}"))
				//case bool:
				// quit out
			case store.Message:
				sender := store.SessionMap[store.Online()[d.(store.Message).Author]].User
				message := &struct {
					ID    uint   `json:"ID"`
					From  string `json:"from"`
					Topic string `json:"topic"`
					Text  string `json:"text"`
				}{ID: d.(store.Message).ID, From: sender.Name, Topic: d.(store.Message).Topic, Text: d.(store.Message).Body}
				msg, _ := json.Marshal(message)
				store.SessionMap[token].WsConn.WriteMessage(1, []byte("{\"type\":\"message\",\"message\":"+string(msg)+"}"))
			}
		}
	}
}
