package chat

import "encoding/json"

//Server : server to add,remove and hold the users
type Server struct {
	Users      map[*User]bool
	Broadcast  chan []byte
	AddUser    chan *User
	DeleteUser chan *User
}

//Start : starts the server.
func (server *Server) Start() {
	for {
		select {
		case conn := <-server.AddUser:
			server.Users[conn] = true
			jsonMsg, _ := json.Marshal(&Message{Content: "/A new user has connected."})
			server.Send(jsonMsg, conn)
		case conn := <-server.DeleteUser:
			if _, ok := server.Users[conn]; ok {
				//closing the channel
				close(conn.Send)
				delete(server.Users, conn)
				jsonMsg, _ := json.Marshal(&Message{Content: "/A user has disconnected."})
				server.Send(jsonMsg, conn)
			}
		case msg := <-server.Broadcast:
			for conn := range server.Users {
				select {
				case conn.Send <- msg:
				default:
					close(conn.Send)
					delete(server.Users, conn)
				}
			}
		}
	}
}

//Send : send msg.
func (server *Server) Send(msg []byte, ignore *User) {
	for conn := range server.Users {
		if conn != ignore {
			conn.Send <- msg
		}
	}
}
