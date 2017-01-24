package main

import (
	"net/http"
	"github.com/gorilla/websocket"
)

func WsNewConn(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "anonymos"
	}
	login := r.URL.Query().Get("login")
	if login == "" {
		login = "anonymos"
	}
	chatName := r.URL.Query().Get("chatname")
	if chatName == "" {
		chatName = "public"
	}
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	usr := &User{
		Login: login,
		Username: username,
	}
	cli := &client{
		wsConn: ws,
		User: usr,
		chat: Chatter,
		send: make(chan []byte, 100),
	}
	chat{}.Connect(chatName, cli, Logger)
}

func WsNewChat(w http.ResponseWriter, r *http.Request) {
	chatName := r.URL.Query().Get("name")
	for name := range AllChats {
		if name == chatName {
			w.Write([]byte("Chat already exists"))
			return
		}
	}
	chat := NewChat(chatName)
	go chat.Run(Logger)
	w.Write([]byte("Chat "+ chatName + " created."))
}
