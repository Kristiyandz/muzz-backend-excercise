package apis

import (
	"database/sql"
	"net/http"
)

/*

	https://api.com/swipe?user_id={my_user_id}&target_user_id={target_userId}&match=YES
	right - NO
	left - YES?
*/

func SwipeHandler(w http.ResponseWriter, r *http.Request) {
	currentUserId := r.Query().Get("user_id")
	targetUserId := r.Query().Get("target_user_id")
	direction := r.Query().Get("direction")

	// logic... DB?
	/*
		user schema -
		user 1 swipes right for user 2
		user one has a record
		user two has a record

		user one record updates swipe_target to be the target ID
		user one record updates swipe_status to be YES

		matches table?
		match_one_id match_two_id user_one_pref user_two_pref

		First swipe will check the user two pref
		 - if no value, insert record
		 - if pref value is NO, return


	*/

	interactionsQuery := `SELECT EXISTS (
		SELECT 1
		FROM user_interactions
		WHERE user_id = UUID_TO_BIN(?),
		AND target_user_id = UUID_TO_BIN(?),
		AND choice = 'YES'
	) AS is_match;`

	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(interactionsQuery, targetUserId, currentUserId)
	if err != nil {
		// handle error...
	}

}
