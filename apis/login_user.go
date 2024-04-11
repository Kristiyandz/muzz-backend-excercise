package apis

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
	generatejwt "github.com/Kristiyandz/muzz-backend-excercise/pkg/generate_jwt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	var reqBody user.UserLoginRequestBody

	errUnmarshal := json.NewDecoder(r.Body).Decode(&reqBody)

	if errUnmarshal != nil {
		http.Error(w, errUnmarshal.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var storedHashedPassword string
	var userUUID string
	query := "SELECT password_hash, uuid FROM users WHERE email = ?"

	err = db.QueryRow(query, reqBody.Email).Scan(&storedHashedPassword, &userUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(reqBody.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	jwt := generatejwt.GenerateJWT(reqBody.Email, userUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := user.UserLoginResponseBody{
		Token: jwt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
