package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/tttmaximttt/go-chat-example/chat"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (self *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.once.Do(func() {
		self.templ = template.Must(template.ParseFiles(filepath.Join("view", self.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	fmt.Println(data)
	self.templ.Execute(w, data)
}

func main() {
	// setup gomniauth
	gomniauth.SetSecurityKey("@ hello @ world @")
	gomniauth.WithProviders(
		//facebook.New("key", "secret",
		//	"http://localhost:8080/auth/callback/facebook"),
		//github.New("key", "secret",
		//	"http://localhost:8080/auth/callback/github"),
		google.New(
			"825934878271-v1qvn7carrogetiqmh83nunrml3f15mo.apps.googleusercontent.com",
			"vS1tGPSxXfCRWMCUz8ffKaBr",
			"http://localhost:8080/auth/callback/google",
		),
	)

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse()

	r := chat.NewRoom(chat.UseFileSystemAvatar)
	//r.trace = trace.New(os.Stdout)

	http.Handle(
		"/assets/",
		http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))),
	)
	http.Handle("/", chat.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/auth/", chat.InitialAuthHandler)
	http.HandleFunc("/uploader", chat.UploaderHandler)

	// LOGOUT
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	// get the room going
	go chat.Run(r)
	// start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	fmt.Printf("Server listening at port %s", addr)
}
