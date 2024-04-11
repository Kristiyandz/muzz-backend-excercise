package apis

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
	hashpassword "github.com/Kristiyandz/muzz-backend-excercise/pkg/hash_password"
	randomchoide "github.com/Kristiyandz/muzz-backend-excercise/pkg/random_choice"
	"github.com/brianvoe/gofakeit/v7"
)

func CreateRandomUserHandler(w http.ResponseWriter, r *http.Request) {
	randomPassword := gofakeit.Password(true, true, true, true, false, 14)
	hashedPassword, err := hashpassword.HashPassword(randomPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	choices := []interface{}{"Male", "Female", "Other", "Prefer not to say"}
	randomUser := user.User{
		Email:    gofakeit.Email(),
		Password: randomPassword,
		Name:     gofakeit.FirstName() + " " + gofakeit.LastName(),
		Gender:   randomchoide.RandomChoice(choices...).(string),
		Age:      gofakeit.Number(18, 100),
	}

	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "INSERT INTO users (email, password_hash, name, gender, age) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?)"
	result, err = db.Exec(query, randomUser.Email, hashedPassword, randomUser.Name, randomUser.Gender, randomUser.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := user.UserCreateResponseBody{
		ID:       int(id),
		Email:    randomUser.Email,
		Password: randomUser.Password,
		Name:     randomUser.Name,
		Gender:   randomUser.Gender,
		Age:      randomUser.Age,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomUser)

}
