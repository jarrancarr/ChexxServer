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
	"github.com/jarrancarr/ChexxServer/match"
	"github.com/jarrancarr/ChexxServer/store"
)

var DEBUG = false

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
			if DEBUG {
				log.Println("...conn read loop")
			}
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
				Move    string `json:"move"`
				Game    string `json:"game"`
			}{}
			read := bytes.NewReader(message)
			err = json.NewDecoder(read).Decode(msg)
			if err != nil {
				log.Println("decode:", err)
				break
			}
			if DEBUG {
				log.Printf("......type: %s  token: %s  Date: %v  Message: %s  Move: %s\n", msg.Type, msg.Token, msg.Epoc, msg.Message, msg.Move)
			}
			if store.SessionMap[msg.Token] != nil {
				switch msg.Type {
				case "ping":
					store.Sessions()[msg.Token].Inbox <- fmt.Sprintf("bounce||%d", time.Now().Unix())
				case "login":
					store.SessionMap[msg.Token].WsConn = conn
					go wsDataQueue(msg.Token)
				case "blitz":
					match.StartBlitz(msg.Token, msg.Game)
				case "abort-blitz":
					match.AbortBlitz(msg.Token)
				case "blitz-move":
					match.BlitzMove(msg.Token, msg.Move)
				case "resign":
					match.BlitzEnd(msg.Token, "Resigned")
				case "blitz-timesup":
					match.BlitzEnd(msg.Token, "Lost on Time")

				}
			} else {
				log.Println("......No token or no session")
			}
		}
	}()
}

func wsDataQueue(token string) {
	if DEBUG {
		log.Println("wsDataQueue")
	}
	if store.SessionMap[token].WsConn != nil {
		live := true
		for live {
			if DEBUG {
				log.Println("...wsDataQueue loop")
			}
			d := <-store.SessionMap[token].Inbox
			if DEBUG {
				log.Printf("......wsDataQueue input")
			}
			switch d.(type) {
			case bool:
				fmt.Printf("...bye bye...")
				live = false
			case string:
				if DEBUG {
					log.Printf("......wsDataQueue string: %v\n", d)
				}
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
				if DEBUG {
					log.Printf("   queue processing match...")
				}
				match, _ := json.Marshal(d)
				if DEBUG {
					log.Printf("1...")
				}
				m := d.(*store.Match)
				white := ""
				black := ""
				if m.BlackPlayerId != 0 {
					b, _ := json.Marshal(store.GetUser(m.BlackPlayerId))
					black = string(b)
				}
				if m.WhitePlayerId != 0 {
					w, _ := json.Marshal(store.GetUser(m.WhitePlayerId))
					white = string(w)
				}
				if DEBUG {
					log.Printf("2...")
				}
				gameType := "view"
				if store.Online()[m.BlackPlayerId] == token || store.Online()[m.WhitePlayerId] == token {
					gameType = "blitz"
				}
				store.SessionMap[token].WsConn.WriteMessage(1, []byte("{\"type\":\""+gameType+"\",\"white\":"+white+",\"black\":"+black+",\"match\":"+string(match)+"}"))
				if DEBUG {
					log.Printf("processed\n")
				}
				//case bool:
				// quit out
			case store.Message:
				sender := store.SessionMap[store.Online()[d.(store.Message).Author]].User
				message := &struct {
					ID     uint   `json:"ID"`
					Author string `json:"userid"`
					Meta   string `json:"meta"`
					From   string `json:"from"`
					Topic  string `json:"topic"`
					Text   string `json:"text"`
				}{ID: d.(store.Message).ID, Author: sender.UserId, Meta: d.(store.Message).Meta, From: sender.Name, Topic: d.(store.Message).Topic, Text: d.(store.Message).Body}
				msg, _ := json.Marshal(message)
				store.SessionMap[token].WsConn.WriteMessage(1, []byte("{\"type\":\"message\",\"message\":"+string(msg)+"}"))
			}
		}
	}
}
