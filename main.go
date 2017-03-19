package main

import (
	"log"
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"fmt"
	//"os"

	//"github.com/tttmaximttt/go-chat-example/trace"
	chat "github.com/tttmaximttt/go-chat-example/chat"
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

	r := chat.NewRoom()
	//r.trace = trace.New(os.Stdout)

	http.Handle(
		"/assets/",
		http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))),
	)
	http.Handle("/", chat.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)

	http.HandleFunc("/auth/", chat.LoginHandler)
	// get the room going
	go chat.Run(r)
	// start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	fmt.Printf("Server listening at port %s", addr)
}