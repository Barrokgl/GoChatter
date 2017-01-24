package main

import (
	"log"
	"net/http"
	"os"
)

var (
	Chatter = NewChat("public")
	Logger = log.New(os.Stdout, "[goChatter]: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Awesome chat server based on go!"))
	})
	http.HandleFunc("/ws/conn", WsNewConn)
	http.HandleFunc("/ws/chat", WsNewChat)
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
	
	// run default public chat
	go Chatter.Run(Logger)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
