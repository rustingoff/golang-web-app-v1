package main

import (
	"net/http"

	"github.com/rustingoff/models"
	"github.com/rustingoff/routes"
	"github.com/rustingoff/utils"
)

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
