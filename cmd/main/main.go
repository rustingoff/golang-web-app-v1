package main

import (
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var client *redis.Client
var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func IndexHandler (w http.ResponseWriter, r *http.Request){
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}

