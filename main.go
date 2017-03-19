package main

import (
	"log"
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"fmt"
	"os"

	"github.com/tttmaximttt/go-chat-example/trace"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (self *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.once.Do(func() {
		self.templ = template.Must(template.ParseFiles(filepath.Join("view", self.filename)))
	})
	self.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse()

	r := newRoom()
	r.trace = trace.New(os.Stdout)

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	fmt.Printf("Server listening at port %s", addr)
}