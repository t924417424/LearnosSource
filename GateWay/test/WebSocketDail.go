package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	log.SetFlags(log.Lshortfile)
	http.HandleFunc("/ws", ws)
	http.HandleFunc("/client", wss)
	http.ListenAndServe(":8011", nil)
}

func ws(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer wsConn.Close()
	for {
		_,msg,err := wsConn.ReadMessage()
		if err != nil {
			log.Fatal(err.Error())
		}
		wsConn.WriteMessage(websocket.TextMessage,msg)
	}
}

func wss(w http.ResponseWriter, r *http.Request) {
	wsConn2, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer wsConn2.Close()
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8011", Path: "/ws"}
	c,_,err := websocket.DefaultDialer.Dial(u.String(),nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer c.Close()
	go func() {
		for {
			_,msg,err := c.ReadMessage()
			if err != nil {
				log.Fatal(err.Error())
			}
			_ = wsConn2.WriteMessage(websocket.TextMessage,msg)
		}
	}()
	for {
		_,msg,err := wsConn2.ReadMessage()
		if err != nil {
			log.Fatal(err.Error())
		}
		_ = c.WriteMessage(websocket.TextMessage,msg)
	}
}
