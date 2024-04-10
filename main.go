package main

import (
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/apis"
)

func main() {
	http.HandleFunc("/user/create", apis.CreateRandomUserHandler)
	http.ListenAndServe(":8080", nil)
}
