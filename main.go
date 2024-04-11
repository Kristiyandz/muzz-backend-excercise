package main

import (
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/apis"
	"github.com/Kristiyandz/muzz-backend-excercise/apis/middleware"
)

func main() {
	http.HandleFunc("/login", apis.LoginUserHandler)

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuthMiddleware(http.HandlerFunc(apis.CreateRandomUserHandler)).ServeHTTP(w, r)
	})

	http.HandleFunc("/discover", func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuthMiddleware(http.HandlerFunc(apis.DiscoverUsersHandler)).ServeHTTP(w, r)
	})
	http.ListenAndServe(":8080", nil)
}
