package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func hellohandler (w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Hello World")
}

func goodbyehandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "GoodBye")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hellohandler).Methods("GET")
	r.HandleFunc("/goodbye", goodbyehandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}