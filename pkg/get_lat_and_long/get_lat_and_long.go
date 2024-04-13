package getlatandlong

import (
	"database/sql"
	"fmt"
)

func GetLatLong(db *sql.DB, userID interface{}) (float64, float64, error) {
	var latitude, longitude float64

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT latitude, longitude FROM users WHERE id = ?")
	if err != nil {
		return 0, 0, err // Return zeros and the error
	}
	defer stmt.Close()

	// Execute the query with the user ID
	err = stmt.QueryRow(userID).Scan(&latitude, &longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no result found
			return 0, 0, fmt.Errorf("no user found with ID %s", userID)
		}
		return 0, 0, err // Return the error encountered
	}

	return latitude, longitude, nil
}
