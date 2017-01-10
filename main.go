package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	} else {
		port = ":" + port
	}
	log.Print("Starting server on port: ", port)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Awesome chat server based on go!")
	})
	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	defer ws.Close()
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		if err = ws.WriteMessage(messageType, p); err != nil {
			panic(err)
		}
		println("Data incoming: ", messageType, p)
	}
}
