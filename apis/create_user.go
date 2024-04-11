package apis

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CreateRandomUserHandler(w http.ResponseWriter, r *http.Request) {
	randomPassword := gofakeit.Password(true, true, true, true, false, 14)
	hashedPassword, err := hashPassword(randomPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	randomUser := user.User{
		Email:    gofakeit.Email(),
		Password: randomPassword,
		Name:     gofakeit.FirstName() + " " + gofakeit.LastName(),
		Gender:   gofakeit.Animal(),
		Age:      gofakeit.Number(18, 100),
	}

	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	uuid := uuid.New().String()
	query := "INSERT INTO users (uuid, email, password_hash, name, gender, age) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, uuid, randomUser.Email, hashedPassword, randomUser.Name, randomUser.Gender, randomUser.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomUser)

}
