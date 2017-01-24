package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	Chatter = NewChat("public")
	Logger = log.New(os.Stdout, "[goChatter]: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Awesome chat server based on go!"))
	})
	http.HandleFunc("/ws/conn", wsNewConn)
	http.HandleFunc("/ws/chat", wsNewChat)
}

func wsNewConn(w http.ResponseWriter, r *http.Request) {
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

func wsNewChat(w http.ResponseWriter, r *http.Request) {
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

func main() {
	// error recover
	defer func() {
		if r := recover(); r != nil {
			Logger.Fatalln(r)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}
	Logger.Print("Starting server on port: ", port)

	go Chatter.Run(Logger)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
