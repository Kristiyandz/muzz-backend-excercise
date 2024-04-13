package apis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	querymapper "github.com/Kristiyandz/muzz-backend-excercise/models/query_mapper"
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
	authenticatedUserID := r.Context().Value("user_id").(float64)
	authUserIdStr := strconv.FormatFloat(authenticatedUserID, 'f', -1, 64)

	fmt.Println(authenticatedUserID)
	// if userID, ok := authenticatedUserID.(int); ok {
	// 	// userID is now a string, and you can use it as such
	// 	fmt.Println("Authenticated User ID:", userID)
	// } else {
	// 	http.Error(w, "Invalid user ID", http.StatusInternalServerError)
	// 	return
	// }

	// query parameter to sort by rank (attractiveness)
	sortBy := r.URL.Query().Get("sort_by")

	// Check if the sortBy parameter is set to "rank"
	// Additional validation could be added here to ensure that the sortBy parameter is only set to "rank"
	var shouldSortByRank bool
	if sortBy == "rank" {
		shouldSortByRank = true
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	loggedInUserDao := querymapper.ExtendedUserDAO(db)

	// Query the database for all users and the authenticated user's latitude and longitude (could be done in a single query with a JOIN, but for simplicity, we'll do two separate queries)
	// query := "SELECT * FROM users WHERE id != ?"
	// loggedInUserLatAndLongQuery := "SELECT latitude, longitude FROM users WHERE id = ?"
	// withRankQuery := `
	// 		SELECT
	// 		u.id AS user_id,
	// 		u.email,
	// 		u.password_hash,
	// 		u.name,
	// 		u.gender,
	// 		u.age,
	// 		u.latitude,
	// 		u.longitude,
	// 		u.created_at,
	// 		u.updated_at,
	// 		IFNULL(COUNT(i.swipe_direction), 0) AS yes_swipes
	// 	FROM users u
	// 	LEFT JOIN interactions i ON u.id = i.swiped_id AND i.swipe_direction = 'YES'
	// 	WHERE u.id != ?
	// 	GROUP BY u.id, u.email, u.password_hash, u.name, u.gender, u.age, u.latitude, u.longitude, u.created_at, u.updated_at
	// 	ORDER BY yes_swipes DESC;
	// 	`

	// Set the query to sort by rank if the sortBy parameter is set
	if shouldSortByRank {
		fmt.Println("Swapping query to sort by rank")
		// query = withRankQuery
	}

	// Query the database for all users
	rows, err := loggedInUserDao.FetchAllUsers(authUserIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the latitude and longitude of the authenticated user
	lat, long, err := getlatandlong.GetLatLong(db, authUserIdStr)
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

		// Perform a different scan based on whether the query is sorting by rank
		if shouldSortByRank {
			err := rows.Scan(&users.ID, &users.Email, &users.Password, &users.Name, &users.Gender, &users.Age, &users.Latitude, &users.Longitude, &users.CreatedAt, &users.UpdatedAt, &users.YesSwipes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := rows.Scan(&users.ID, &users.Email, &users.Password, &users.Name, &users.Gender, &users.Age, &users.Latitude, &users.Longitude, &users.CreatedAt, &users.UpdatedAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
