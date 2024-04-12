package apis

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
)

func DiscoverUsersHandler(w http.ResponseWriter, r *http.Request) {

	authenticatedUserID := r.Context().Value("user_id")

	query := "SELECT * FROM users"

	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users user.UsersTableRecord
	var allUsers []user.DiscoverUserResponseBody
	for rows.Next() {
		// need to be SCANNED IN A ROW
		err := rows.Scan(&users.ID, &users.Email, &users.Password, &users.Name, &users.Gender, &users.Age, &users.CreatedAt, &users.UpdatedAt, &users.Latitude, &users.Longitude)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allUsers = append(allUsers, user.DiscoverUserResponseBody{
			ID:        users.ID,
			Name:      users.Name,
			Gender:    users.Gender,
			Age:       users.Age,
			Latitude:  users.Latitude,
			Longitude: users.Longitude,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUsers)

}
