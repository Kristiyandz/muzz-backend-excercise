package apis

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
	hashpassword "github.com/Kristiyandz/muzz-backend-excercise/pkg/hash_password"
	randomchoice "github.com/Kristiyandz/muzz-backend-excercise/pkg/random_choice"
	"github.com/brianvoe/gofakeit/v7"
)

func CreateRandomUserHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a random password & hash it
	randomPassword := gofakeit.Password(true, true, true, true, false, 14)
	hashedPassword, err := hashpassword.HashPassword(randomPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a random gender selection
	choices := []interface{}{"Male", "Female", "Other", "Prefer not to say"}
	randomGenderChoice := randomchoice.RandomChoiceFromSlice(choices).(string)

	// Create a random user
	randomUser := user.User{
		Email:     gofakeit.Email(),
		Password:  randomPassword,
		Name:      gofakeit.FirstName() + " " + gofakeit.LastName(),
		Gender:    randomGenderChoice,
		Age:       gofakeit.Number(18, 100),
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}

	// Insert the random user into the database
	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, "/user/create cannot connect to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "INSERT INTO users (email, password_hash, name, gender, age, latitude, longitude) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, randomUser.Email, hashedPassword, randomUser.Name, randomUser.Gender, randomUser.Age, randomUser.Latitude, randomUser.Longitude)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the ID of the inserted user
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "/user/create cannot get last iserted index", http.StatusInternalServerError)
		return
	}

	// Return the created user
	response := user.UserCreateResponseBody{
		ID:       int(id),
		Email:    randomUser.Email,
		Password: randomUser.Password,
		Name:     randomUser.Name,
		Gender:   randomUser.Gender,
		Age:      randomUser.Age,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
