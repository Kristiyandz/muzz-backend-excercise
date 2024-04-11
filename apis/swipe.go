package apis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	fmt.Println("currentUserId", currentUserId)
	fmt.Println("targetUserId", targetUserId)
	fmt.Println("match", match)

	interactionsQuery := `SELECT EXISTS (
		SELECT 1
		FROM interactions
		WHERE user_id = ?
		AND target_user_id = ?
		AND choice = ?
	) AS is_match;`

	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/muzzmaindb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(interactionsQuery, targetUserId, currentUserId, match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var isMatch bool
	for rows.Next() {
		err := rows.Scan(&isMatch)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if isMatch {
		matchQuery := `INSERT INTO matches(user1_id, user2_id, created_at) VALUES (?, ?, ?)`
		_, err := db.Exec(matchQuery, currentUserId, targetUserId, time.Now())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userIdIntValue, err := strconv.Atoi(targetUserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result := user.MatchedResultResponseBody{
			Match:   true,
			MatchID: userIdIntValue,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

	} else {
		interactionsInsertQuery := `INSERT INTO interactions (user_id, target_user_id, choice, created_at) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(interactionsInsertQuery, currentUserId, targetUserId, match, time.Now())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
