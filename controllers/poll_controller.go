package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"go-api/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func Poll_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not", http.StatusMethodNotAllowed)
		return
	}

	pollId := r.PathValue("pollId")
	if pollId == "" {
		http.Error(w, "Missing pollId parameter", http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(pollId)
	if err != nil {
		http.Error(w, "Invalid pollID parameter", http.StatusBadRequest)
		return
	}


	db, err := models.ConnectDB()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT id, title, createdat, updatedat FROM poll WHERE id = $1", pollId)

	var poll models.Poll
	err = row.Scan(&poll.Id, &poll.Title, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Poll Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch poll", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	rows, err := db.Query("SELECT id, title, pollid FROM polloption WHERE pollid = $1", pollId)
	if err != nil {
		http.Error(w, "Failed to fetch poll options", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	type Option struct {
		Id    string `json:"id"`
		Title string `json:"title"`
	}

	var pollOptions []Option

	for rows.Next() {
		var option models.PollOption
		err := rows.Scan(&option.Id, &option.Title, &option.PollId)
		if err != nil {
			http.Error(w, "Failed to fetch poll options", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		var pollOptoinCleaned Option
		pollOptoinCleaned.Id = option.Id
		pollOptoinCleaned.Title = option.Title

		pollOptions = append(pollOptions, pollOptoinCleaned)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to fetch poll options", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	response := struct {
		Poll        models.Poll `json:"poll"`
		PollOptions []Option    `json:"pollOptions"`
	}{
		Poll:        poll,
		PollOptions: pollOptions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
