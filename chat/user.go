package chat

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type User struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

func (c *User) Read(server *Server) {
	defer func() {
		server.DeleteUser <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			server.DeleteUser <- c
			c.Socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: string(message)})
		server.Broadcast <- jsonMessage
	}
}

func (c *User) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
