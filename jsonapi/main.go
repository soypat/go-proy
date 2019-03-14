package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)
func main() {
	err := main2()
	if err != nil {
		log.Fatal("Error en programa main: ",err)
	} else {
		log.Fatal("Server terminado.")
	}

}
var store *Store

func helloWorld(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello World\n"))
}

type Server struct {
	r *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.r.ServeHTTP(w, r)
}
func ListTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tasks, err := store.GetTasks()
	if err != nil {
		log.Fatal("Error storing acquiring tasks from db during listing: ", err)
	}
	b, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal("Error Marshalling json Tasks during listing: ", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func CreateTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := Task{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if t.Id == 0 && t.Title == "" {
		w.WriteHeader(http.StatusBadRequest) // TODO this part of code is an eyesore
		return
	}
	if err != nil {
		log.Fatal("Error decoding json during Task creation: ", err)
	}
	err = store.CreateTask(&t)
	if err != nil {
		log.Fatal("Error storing task during task creation: ", err)
	}
	w.WriteHeader(http.StatusCreated)
}

func LoginTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := Task{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&t)
	fmt.Printf("%d , %s", t.Id, t.Title)
	if t.Id == 0 && t.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		//return
		log.Fatal("Error decoding json in Login: ", err)
	}
	err = store.CreateTask(&t)
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		//return
		log.Fatal("Error storing task to db: ", err)
	}
	w.WriteHeader(http.StatusCreated)
}

func ReadTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "ReadTask\n")
}

func UpdateTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "UpdateTask\n")
}

func DeleteTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "DeleteTask\n")
}

func main2() error {

	router := httprouter.New()
	router.GET("/", ListTasks)

	router.POST("/login", LoginTask)
	router.POST("/", CreateTask)

	router.GET("/:id", UpdateTask)
	router.PUT("/:id", UpdateTask)
	router.DELETE("/:id", DeleteTask)

	var err error
	store, err = NewStore()

	if err != nil {

		return fmt.Errorf(fmt.Sprintf("Error guardando a db: ",err) )
	}

	store.Initialize()
	defer store.Close()

	err = http.ListenAndServe(":8080", &Server{router})
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error en Listen/Serve: ",err) )
	}
	return nil
}
