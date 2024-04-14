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
		`SELECT EXISTS (
		    SELECT 1
		    FROM interactions as a
		    JOIN interactions as b ON a.swiper_id = b.swiped_id AND a.swiped_id = b.swiper_id
		    WHERE a.swiper_id = ? AND a.swiped_id = ? 
		      AND a.swipe_direction = 'YES' AND b.swipe_direction = 'YES'
		) AS is_match;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(currentUserId, targetUserId)
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
