package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", StaticHandler)
	http.ListenAndServe(":8080", nil)
}

type User struct {
	Name string
	Age  int
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		Name: "example",
		Age:  20,
	}
	tmpl := template.Must(template.ParseFiles("./views/index.tmpl"))
	tmpl.Execute(w, user)
}