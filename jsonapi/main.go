package main


import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
)
func main() {
	err := main2()
	if err != nil {
		log.Fatal("Error en programa main: ",err)
	} else {
		log.Println("Server terminado.")
	}

}
var store *Store

func IPget() ([]byte,error) {
	response, err := http.Get("https://ifconfig.co/")
	if err != nil {
		return nil,err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return body,err
	}
	return body,nil
}

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

func LandPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	bytes,err :=ioutil.ReadFile("landPage.html")
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError) // TODO this part of code is an eyesore
		return
	}
	w.Write(bytes)

	w.WriteHeader(http.StatusCreated)
}

func getStyle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	bytes,err :=ioutil.ReadFile("CSS/basic.css")
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError) // TODO this part of code is an eyesore
		return
	}
	w.Write(bytes)

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
	//client := &http.Client{Timeout:8*time.Second,}
	//fmt.Println("Acquiring IP...")

	myIP,err := IPget()

	if err != nil {
		fmt.Println("Failed IP Get")
	} else {
		fmt.Printf("Got IP: %s\nStarting Server...\n\n",string(myIP))
	}

	router := httprouter.New()
	router.GET("/", LandPage)
	router.GET("/tasks", ListTasks)
	router.POST("/", CreateTask)
	router.GET("/CSS/basic.css", getStyle)

	router.POST("/login", LoginTask)



	//router.GET("/:id", ReadTask)
	router.PUT("/:id", UpdateTask)
	router.DELETE("/:id", DeleteTask)

	store, err = NewStore()

	if err != nil {

		return fmt.Errorf(fmt.Sprintf("Error guardando a db: ",err) )
	}

	store.Initialize()
	defer store.Close()

	err = http.ListenAndServe(":80", &Server{router})
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error en Listen/Serve: ",err) )
	}
	return nil
}
