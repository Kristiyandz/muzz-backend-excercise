package apis

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	querymapper "github.com/Kristiyandz/muzz-backend-excercise/models/query_mapper"
	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
)

func SwipeHandler(w http.ResponseWriter, r *http.Request) {
	currentUserId := r.URL.Query().Get("user_id")
	targetUserId := r.URL.Query().Get("target_user_id")
	match := r.URL.Query().Get("match")

	/*
		Ideally all of the required parameters should be validated here or in a middleware
		to ensure that they are not empty and are of the correct type.
		For the sake of simplicity, we will only check if they are empty.
	*/

	if currentUserId == "" || targetUserId == "" || match == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	if currentUserId == targetUserId {
		http.Error(w, "You cannot swipe on yourself", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, "/swpie cannot connect to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	lggedinUserDao := querymapper.ExtendedUserSwipeDAO(db)

	rows, err := lggedinUserDao.CheckUserInteractions(currentUserId, targetUserId)
	if err != nil {
		http.Error(w, "/swipe cannot query db", http.StatusInternalServerError)
		return
	}

	rankingRows, err := lggedinUserDao.ApplyRanking(currentUserId)
	if err != nil {
		http.Error(w, "/swipe cannot query db", http.StatusInternalServerError)
		return
	}

	var userId, yesSwipes int
	for rankingRows.Next() {
		err := rankingRows.Scan(&userId, &yesSwipes)
		if err != nil {
			http.Error(w, "/swipe cannot scan ranking rows", http.StatusInternalServerError)
			return
		}
	}

	if err := rankingRows.Err(); err != nil {
		http.Error(w, "/swipe cannot iterate ranking rows", http.StatusInternalServerError)
		return
	}

	// Check if the user has already swiped on the target user
	var isMatch bool
	for rows.Next() {
		err := rows.Scan(&isMatch)
		if err != nil {
			http.Error(w, "/swipe cannot scan rows", http.StatusInternalServerError)
			return
		}
	}

	// If the user has already swiped on the target user and the choice is "match", create a match
	if isMatch {
		matchQuery := `INSERT INTO matches(user1_id, user2_id, created_at) VALUES (?, ?, ?)`
		_, err := db.Exec(matchQuery, currentUserId, targetUserId, time.Now())
		if err != nil {
			http.Error(w, "/swipe cannot execute match query", http.StatusInternalServerError)
			return
		}

		userIdIntValue, err := strconv.Atoi(targetUserId)
		if err != nil {
			http.Error(w, "/swipe failed to convert str to int", http.StatusInternalServerError)
			return
		}

		// Return the match response
		matchResult := user.MatchedResultResponseBody{
			Match:   true,
			MatchID: &userIdIntValue,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matchResult)

	} else {
		// If the user has not already swiped on the target user, insert the swipe into the interactions table
		interactionsInsertQuery := `INSERT INTO interactions (swiper_id, swiped_id, swipe_direction, created_at) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(interactionsInsertQuery, currentUserId, targetUserId, match, time.Now())
		if err != nil {
			http.Error(w, "/swipe failed to insert into interactions", http.StatusInternalServerError)
			return
		}

		// Return the swipe result
		swipeResult := user.MatchedResultResponseBody{
			Match:   false,
			MatchID: nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(swipeResult)
	}

}
