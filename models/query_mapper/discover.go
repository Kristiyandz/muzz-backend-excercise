package querymapper

import (
	"database/sql"
	"fmt"
)

type UserDAO struct {
	db *sql.DB
}

func ExtendedUserDAO(db *sql.DB) UserDAO {
	return UserDAO{db: db}
}

func (dao *UserDAO) FetchAllUsers(loggedInUserID string) (*sql.Rows, error) {
	stmt, err := dao.db.Prepare(`
		SELECT *
			FROM users
			WHERE id != ?
				AND id NOT IN (
					SELECT swiped_id
					FROM interactions
					WHERE swipe_direction = 'YES'
				)
		`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(loggedInUserID)
}

func (dao *UserDAO) FetchLoggedUserLatLong(loggedInUserID string) (*sql.Rows, error) {
	stmt, err := dao.db.Prepare("SELECT latitude, longitude FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(loggedInUserID)
}

func (dao *UserDAO) FetchUsersWithRank(loggedInUserID string) (*sql.Rows, error) {
	stmt, err := dao.db.Prepare(
		`SELECT
			u.id AS user_id,
			u.email,
			u.password_hash,
			u.name,
			u.gender,
			u.age,
			u.latitude,
			u.longitude,
			u.created_at,
			u.updated_at,
			IFNULL(COUNT(i.swipe_direction), 0) AS yes_swipes
		FROM users u
		LEFT JOIN interactions i ON u.id = i.swiped_id AND i.swipe_direction = 'YES'
		WHERE u.id != ?
		GROUP BY u.id, u.email, u.password_hash, u.name, u.gender, u.age, u.latitude, u.longitude, u.created_at, u.updated_at
		ORDER BY yes_swipes DESC;`)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(loggedInUserID)
}

func (dao *UserDAO) SortByAgeOrGender(loggedInUserID, sortBy string) (*sql.Rows, error) {
	query := "SELECT * FROM users WHERE id != ?"

	// Append sorting to the query
	if sortBy == "age" || sortBy == "gender" {
		query += fmt.Sprintf(" ORDER BY %s ASC", sortBy)
	}

	stmt, err := dao.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(loggedInUserID)
}
