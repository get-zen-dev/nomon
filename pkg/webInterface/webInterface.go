package webInterface

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	staticFileDirectory := http.Dir("./")
	staticFileHandler := http.StripPrefix("./", http.FileServer(staticFileDirectory))
	r.PathPrefix("./").Handler(staticFileHandler).Methods("GET")
	r.HandleFunc("/breakpoints", createBreakpointsHandler).Methods("POST")

	return r
}

var tpl = template.Must(template.ParseFiles("./index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

type Breakpoints struct {
	X string `json:"x"`
	Y string `json:"y"`
}

func createBreakpointsHandler(w http.ResponseWriter, r *http.Request) {
	b := Breakpoints{}

	err := r.ParseForm()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b.X = r.Form.Get("x")
	b.Y = r.Form.Get("y")

	log.Println(b)
	http.Redirect(w, r, "/", http.StatusFound)
}
