package main

import (
	"fmt"
	"gochat-websocket-server/chat"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var server = &chat.Server{
	Broadcast:  make(chan []byte),
	AddUser:    make(chan *chat.User),
	DeleteUser: make(chan *chat.User),
	Users:      make(map[*chat.User]bool),
}

func main() {
	fmt.Println("Starting chat app server...")
	go server.Start()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":8080", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	user := &chat.User{
		ID:     strconv.Itoa(rand.Intn(100)),
		Socket: conn,
		Send:   make(chan []byte),
	}

	server.AddUser <- user

	go user.Read(server)
	go user.Write()
}
