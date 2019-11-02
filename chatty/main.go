package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

const (secret = "C3q7eh09waMNdr8n89Z8nMPo"
		clientID = "722589219392-r9ff97evqicjecv1h1uro7k5s47gsi2a.apps.googleusercontent.com"
		key = "iop dots like krop nogs"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

// ServeHTTP handles the HTTP request.
// We will define a new type that will take a
//filename string, compile the template once (using the sync.Once type), keep the
//reference to the compiled template, and then respond to HTTP requests. this will be templateHandler
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	// set up gomniauth
	gomniauth.SetSecurityKey(key)
	// , github.New("key", "secret","http://localhost:8080/auth/callback/github")
	gomniauth.WithProviders(google.New(clientID, secret,"http://localhost:8080/auth/callback/google"))

	r := newRoom()
	http.Handle("/", &templateHandler{filename: "lz.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chatjs.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room",r)
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

