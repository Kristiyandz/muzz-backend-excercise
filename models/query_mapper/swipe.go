package querymapper

import (
	"database/sql"
)

type UserSwipeDAO struct {
	db *sql.DB
}

func ExtendedUserSwipeDAO(db *sql.DB) UserSwipeDAO {
	return UserSwipeDAO{db: db}
}

func (dao *UserSwipeDAO) CheckUserInteractions(currentUserId, targetUserId string) (*sql.Rows, error) {
	stmt, err := dao.db.Prepare(
		`SELECT EXISTS (SELECT 1
			FROM   interactions
			WHERE  swiper_id = ?
						 AND swiped_id = ?
						 AND swipe_direction = 'YES') AS is_match; `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(currentUserId, targetUserId)
}

func (dao *UserSwipeDAO) ApplyRanking(loggedInUserID string) (*sql.Rows, error) {
	stmt, err := dao.db.Prepare(
		`SELECT EXISTS (SELECT 1
			FROM   interactions
			WHERE  swiper_id = ?
						 AND swiped_id = ?
						 AND swipe_direction = 'YES') AS is_match; `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(loggedInUserID)
}
