package main

import (
	"encoding/json"
)

type Message struct {
	Message string `json:"message"`
	Sender string `json:"sender"`
	Created string `json:"created"`
}

func FromJson(jsonIn []byte) (message *Message) {
	json.Unmarshal(jsonIn, &message)
	return
}


