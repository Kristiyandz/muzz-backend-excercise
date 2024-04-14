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

func (dao *UserSwipeDAO) CheckUserInteractions(targetUserId, currentUserId string) (bool, error) {

	/*
		@TODO:(BUG): Currently if the first user swipes YES,
		eventually the second user swipes NO,
		a match will still be created.

		Works when the first user swipes NO and the second user swipes YES.
	*/
	stmt, err := dao.db.Prepare(
		`SELECT EXISTS (
			SELECT 1
			FROM interactions
			WHERE swiper_id = ?
			AND swiped_id = ?
			AND swipe_direction = 'YES'
		) AS is_match;`)

	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var isMatch bool
	err = stmt.QueryRow(targetUserId, currentUserId).Scan(&isMatch)
	if err != nil {
		return false, err

	}
	return isMatch, nil
}

func (dao *UserSwipeDAO) ApplyRanking() (*sql.Rows, error) {
	stmt, err := dao.db.Prepare(
		`SELECT swiped_id AS target_user_id,
				Count(*)  AS yes_swipes
			FROM   interactions
			WHERE  swipe_direction = 'YES'
			GROUP  BY swiped_id
			ORDER  BY yes_swipes DESC;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query()
}
