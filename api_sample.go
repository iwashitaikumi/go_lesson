package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/api", ApiHandler)

	http.ListenAndServe(":8080", nil)
}

type User struct {
	Name string
	Age  int
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		Name: "example",
		Age:  20,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}