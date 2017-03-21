package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tttmaximttt/go-chat-example/trace"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be using for broadcasting.
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	trace   trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		trace:   trace.Off(),
	}
}

func (self *room) run() {
	for {
		select {
		case client := <-self.join:
			self.clients[client] = true
			self.trace.Trace("New client joined")
		case client := <-self.leave:
			delete(self.clients, client)
			close(client.send)
			self.trace.Trace("Client leave room")
		case msg := <-self.forward:
			// forward message to all clients
			self.trace.Trace("Message received: ", string(msg))
			for client := range self.clients {
				client.send <- msg
				self.trace.Trace(" -- sent to client")
			}
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func Run(self *room) {
	self.run()
}

func NewRoom() *room {
	return newRoom()
}

func (self *room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := r.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan *message, messageBufferSize),
		room:   self,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	self.join <- client
	defer func() { self.leave <- client }()
	go client.write()
	client.read()
}
