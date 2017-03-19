package chat

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send chan []byte
	room *room
}

func (self *client) read() {
	defer self.socket.Close()

	for{
		_, msg, err := self.socket.ReadMessage()

		if err != nil {
			return
		}

		self.room.forward <- msg
	}
}

func (self *client) write() {
	defer self.socket.Close()
	for msg := range self.send {
		err := self.socket.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			return
		}
	}
}