package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	Chatter = NewChat("gopher")
	Logger = log.New(os.Stdout, "[goChatter]: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Awesome chat server based on go!"))
	})
	http.HandleFunc("/ws", wsHandler)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	usr := &User{
		Login:"userlogin",
		Username:"username",
	}
	cli := &client{
		wsConn: ws,
		User: usr,
		chat: Chatter,
		send: make(chan []byte, 100),
	}
	cli.chat.join <- cli
	go cli.read(Logger)
	go cli.write(Logger)
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
