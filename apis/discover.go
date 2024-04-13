package apis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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

	// query parameter to sort by rank (attractiveness)
	sortBy := r.URL.Query().Get("sort_by")

	// Check if the sortBy parameter is set to "rank"
	// Additional validation could be added here to ensure that the sortBy parameter is only set to "rank"
	var shouldSortByRank, sortByAge, sortByGender bool

	fmt.Println(sortByGender)

	switch sortBy {
	case "rank":
		shouldSortByRank = true
	case "age":
		sortByAge = true
	case "gender":
		sortByGender = true
	}

	// Connect to the database
	// DB connection name and password should be stored in a secure location but for the sake of simplicity, we will hardcode it here
	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	loggedInUserDao := querymapper.ExtendedUserDAO(db)

	var rows *sql.Rows
	// Fetch all users or users sorted by rank(attractiveness)
	if shouldSortByRank {
		rows, err = loggedInUserDao.FetchUsersWithRank(authUserIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if sortByAge {
		rows, err = loggedInUserDao.SortByAgeOrGender(authUserIdStr, "age")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Fetch all users(non-sorted)
		rows, err = loggedInUserDao.FetchAllUsers(authUserIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
			// Perform default scan
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

	// Sort the users by distance from the authenticated user in ascending order (closest first)
	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].DistanceFromMe < allUsers[j].DistanceFromMe
	})

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUsers)

}
