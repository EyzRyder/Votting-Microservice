package util

import (
	"context"
	"database/sql"
	"net/http"

	"go-api/models"

	_ "github.com/lib/pq"
)


func GetPoll( w http.ResponseWriter,pollId string) (*models.Response, error) {


	redisClient, err := models.InitRedis()

	defer redisClient.Close()

	ctx := context.Background()
	if err != nil {
		http.Error(w, "Failed to connect to Redis", http.StatusInternalServerError)
		return nil,err
	}

    rangeResults, err := redisClient.ZRangeWithScores(ctx,pollId,0,-1).Result()
	if err != nil {
		http.Error(w, "Failed to fetch poll options from Redis", http.StatusInternalServerError)
		return nil, err
	}

    pollScores := make(map[string]float64)

    for _,result := range rangeResults{
        pollScores[result.Member.(string)] = result.Score
    }

	db, err := models.ConnectDB()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow("SELECT id, title, createdat, updatedat FROM poll WHERE id = $1", pollId)

	var poll models.Poll
	err = row.Scan(&poll.Id, &poll.Title, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Poll Not Found", http.StatusNotFound)
		    return nil, err
		}
		http.Error(w, "Failed to fetch poll", http.StatusInternalServerError)
		return nil, err
	}

	rows, err := db.Query("SELECT id, title, pollid FROM polloption WHERE pollid = $1", pollId)
	if err != nil {
		http.Error(w, "Failed to fetch poll options", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()


	var pollOptions []models.Option

	for rows.Next() {
		var option models.PollOption
		err := rows.Scan(&option.Id, &option.Title, &option.PollId)
		if err != nil {
			http.Error(w, "Failed to fetch poll options", http.StatusInternalServerError)
		return nil, err
		}
		var pollOptoinCleaned models.Option
		pollOptoinCleaned.Id = option.Id
		pollOptoinCleaned.Title = option.Title
        pollOptoinCleaned.Score = pollScores[option.Id]

		pollOptions = append(pollOptions, pollOptoinCleaned)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to fetch poll options", http.StatusInternalServerError)
		return nil, err
	}

	response := models.Response{
		Poll:        poll,
		PollOptions: pollOptions,
	}

    return &response, nil

}
