package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// helloWorld Two url available: "/" and "/bye"
func helloWorld(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		// en el URL "/" nos fijamos si recibimos un GET o POST request
		switch r.Method {
		// Se puede acceder a GET con cURL de la forma: curl -si "http://localhost:8000/?foo=1&bar=2" creando asi dos keys y sus values
		case "GET":
			// donde k:Key y v:Value
			for k, v := range r.URL.Query() {
				fmt.Printf("%s: %s\n", k, v)
			}
			w.Write([]byte("Received a GET request\n"))
			// mandar POST con: curl -si -X POST -d "data del POST" http://localhost:8000
		case "POST":
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
				fmt.Println("Error reading POST...")
			}
			fmt.Printf("%s\n", reqBody)
			w.Write([]byte("Received a POST request\n"))
		default:
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte(http.StatusText(http.StatusNotImplemented) + "\n"))
		}
	case "/bye":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"goodbye": "cruel world"}`))
	default:
		http.NotFound(w, r)
		return
	}
}

func main() {
	http.HandleFunc("/", helloWorld)
	http.ListenAndServe(":8000", nil)
}
