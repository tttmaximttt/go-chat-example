package chat

import (
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (self *client) read() {
	defer self.socket.Close()

	for {
		var msg *message
		err := self.socket.ReadJSON(&msg)

		if err != nil {
			return
		}

		msg.When = time.Now()
		msg.Name = self.userData["name"].(string)

		if avatarURL, ok := self.userData["avatar"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}

		self.room.forward <- msg
	}
}

func (self *client) write() {
	defer self.socket.Close()
	for msg := range self.send {
		err := self.socket.WriteJSON(msg)

		if err != nil {
			return
		}
	}
}
