package main

import (
	"math/rand"
	"strconv"
)

var (
	AllChats   = make(map[string]*chat)
	chatCounts int
)

type chat struct {
	name      string
	forward chan []byte
	userConns map[*client]bool
	join      chan *client
	leave     chan *client
}

func (c *chat) Run() {
	for {
		select {
		case ch := <-c.forward:
			for user := range c.userConns {
				select {
				case user.send <- ch:
				default:
					delete(c.userConns, user)
					close(user.send)
				}
			}
		case ch := <-c.join:
			Logger.Println("user join: ", ch.Username)
			c.userConns[ch] = true
		case ch := <-c.leave:
			Logger.Println("user leave: ", ch.Username)
			delete(c.userConns, ch)
			close(ch.send)
		}
	}
}

func NewChat(name string) *chat {
	if name == "" {
		num := strconv.Itoa(rand.Int())
		name = "Chat" + num
	}
	chat := &chat{
		name:      name,
		forward: make(chan []byte, 100),
		userConns: make(map[*client]bool),
		join:      make(chan *client, 100),
		leave:     make(chan *client, 100),
	}
	chatCounts++
	AllChats[name] = chat
	Logger.Println("Creating chat: ", name)
	return chat
}

