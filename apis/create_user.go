package apis

import (
	"encoding/json"
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
)

func CreateRandomUserHandler(w http.ResponseWriter, r *http.Request) {
	user := user.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "Test",
		Gender:   "Male",
		Age:      20,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}
