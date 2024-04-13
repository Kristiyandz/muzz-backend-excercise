package querymapper

import (
	"database/sql"
)

type UserDAO struct {
	db *sql.DB
}

func ExtendedUserDAO(db *sql.DB) UserDAO {
	return UserDAO{db: db}
}

func (dao *UserDAO) FetchAllUsers(loggedInUserID string) (*sql.Rows, error) {
	stmt, err := dao.db.Prepare("SELECT * FROM users WHERE id != ?")
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
