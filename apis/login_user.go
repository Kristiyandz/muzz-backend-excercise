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

	// Unmarshal the request body
	var reqBody user.UserLoginRequestBody
	errUnmarshal := json.NewDecoder(r.Body).Decode(&reqBody)

	if errUnmarshal != nil {
		http.Error(w, errUnmarshal.Error(), http.StatusBadRequest)
		return
	}

	// Connect to the database
	// For simplicity, we will hardcode the database connection details here but in a production application, these should be stored securely
	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the database for the user's hashed password
	var storedHashedPassword string
	var userId int
	query := "SELECT password_hash, id FROM users WHERE email = ?"

	// Check if the user exists in the database
	err = db.QueryRow(query, reqBody.Email).Scan(&storedHashedPassword, &userId)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password with the password from the request body
	if err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(reqBody.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token and send it in the response
	jwt := generatejwt.GenerateJWT(reqBody.Email, userId)

	response := user.UserLoginResponseBody{
		Token: jwt,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
