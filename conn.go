package main

import (
	"github.com/gorilla/websocket"
)

type User struct {
	Login    string `json:"login"`
	Username string `json:"username"`
}

type client struct {
	wsConn *websocket.Conn
	*User
	*chat
	send chan []byte
}

func (c *client) read() {
	for {
		if msgT, msg, err := c.wsConn.ReadMessage(); err == nil {
			Logger.Println("reading msg: ", string(msg)," " ,msgT)
			c.chat.forward <- msg
		} else {
			break
		}
	}
	c.wsConn.Close()
}

func (c *client) write() {
	for msg := range c.send {
		Logger.Println("writing msg: ", string(msg))
		if err := c.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.wsConn.Close()
}