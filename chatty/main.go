package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
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
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename))) // filepath joins, wow!
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chatjs.html"})
	http.Handle("/room",r)
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
