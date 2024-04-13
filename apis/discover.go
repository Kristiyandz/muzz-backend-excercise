package apis

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
	calculatedistance "github.com/Kristiyandz/muzz-backend-excercise/pkg/calculate_distance"
	getlatandlong "github.com/Kristiyandz/muzz-backend-excercise/pkg/get_lat_and_long"
)

type coordinate struct {
	lat float64
	lng float64
}

func DiscoverUsersHandler(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user's ID from the context (set by the JWTAuthMiddleware)
	authenticatedUserID := r.Context().Value("user_id")

	// Query the database for all users and the authenticated user's latitude and longitude (could be done in a single query with a JOIN, but for simplicity, we'll do two separate queries)
	query := "SELECT * FROM users WHERE id != ?"
	loggedInUserLatAndLongQuery := "SELECT latitude, longitude FROM users WHERE id = ?"

	// Connect to the database
	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the database for all users
	rows, err := db.Query(query, authenticatedUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	/* Query the database for the authenticated user's latitude and longitude */
	/* The user lat/long can be passes as claims from the JWT to avoid the query but for simplicyty and making it work, the query will do */
	loggedInUserLatLong, err := db.Query(loggedInUserLatAndLongQuery, authenticatedUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer loggedInUserLatLong.Close()

	// Get the latitude and longitude of the authenticated user
	lat, long, err := getlatandlong.GetLatLong(db, authenticatedUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the authenticated user's coordinates
	loggedInUserCoordinates := coordinate{lat, long}

	// Loop through the users and calculate the distance from the authenticated user
	var users user.UsersTableRecord
	var allUsers []user.DiscoverUserResponseBody
	for rows.Next() {
		err := rows.Scan(&users.ID, &users.Email, &users.Password, &users.Name, &users.Gender, &users.Age, &users.Latitude, &users.Longitude, &users.CreatedAt, &users.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Store the requested user's coordinates
		requestedUserCoordinates := coordinate{users.Latitude, users.Longitude}
		// Calculate the distance from the authenticated user (default unit is in miles)
		distanceFromMe := calculatedistance.Distance(loggedInUserCoordinates.lat, loggedInUserCoordinates.lng, requestedUserCoordinates.lat, requestedUserCoordinates.lng)

		// Append the user to the response
		allUsers = append(allUsers, user.DiscoverUserResponseBody{
			ID:             users.ID,
			Name:           users.Name,
			Gender:         users.Gender,
			Age:            users.Age,
			DistanceFromMe: distanceFromMe,
		})
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUsers)

}
