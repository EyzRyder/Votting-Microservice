package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go-api/models"
	"go-api/util"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func VotePoll_Controller(vote *util.VotingPubSub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
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

		var pollOption struct {
			PollOptionId string `json:"pollOptionId"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&pollOption); err != nil {
			http.Error(w, "Invalid requst body", http.StatusBadRequest)
		}

		if pollOption.PollOptionId == "" {
			http.Error(w, "pollOptionId cannot be empty", http.StatusBadRequest)
			return
		}
		_, err = uuid.Parse(pollOption.PollOptionId)
		if err != nil {
			http.Error(w, "Invalid Poll Option Id", http.StatusBadRequest)
			return
		}

		var sessionId string

		cookie, err := r.Cookie("sessionId")

		db, err := models.ConnectDB()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		defer db.Close()

		redisClient, err := models.InitRedis()

		defer redisClient.Close()

		ctx := context.Background()
		if err != nil {
			http.Error(w, "Failed to connect to Redis", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		if err == http.ErrNoCookie || cookie == nil || cookie.Value == "" {
			sessionId = uuid.New().String()
			cookie = &http.Cookie{
				Name:     "sessionId",
				Value:    sessionId,
				Path:     "/",
				MaxAge:   60 * 60 * 24 * 30,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			}
			http.SetCookie(w, cookie)
		} else {
			sessionId = cookie.Value

			row := db.QueryRow("SELECT id, sessionId, pollId, pollOptionId, createdat FROM vote WHERE pollId = $1 AND sessionId = $2",
				pollId, sessionId)

			var vote models.Vote
			err = row.Scan(&vote.Id, &vote.SessionId, &vote.PollId, &vote.PollOptionId, &vote.CreatedAt)

			if err == sql.ErrNoRows {
			} else if err != nil {
				http.Error(w, "Failed to fetch vote", http.StatusInternalServerError)
				log.Fatal(err)
				return
			} else {
				if vote.PollId == pollId && vote.PollOptionId != pollOption.PollOptionId {

					_, err := db.Exec("DELETE FROM vote WHERE pollId = $1 AND sessionId=$2",
						pollId, sessionId)

					if err != nil {
						http.Error(w, "Failed to deletd vote row", http.StatusInternalServerError)
						return
					}

					_, err = redisClient.ZIncrBy(ctx, pollId, -1.0, vote.PollOptionId).Result()
					if err != nil {
						http.Error(w, "Failed to Decrement vote count in Redis", http.StatusInternalServerError)
						log.Fatal(err)
						return
					}

				} else if vote.PollId == pollId && vote.PollOptionId == pollOption.PollOptionId {
					http.Error(w, "You have already voted on this poll", http.StatusBadRequest)
					return
				}
			}

		}

		_, err = db.Exec("INSERT INTO vote (sessionId, pollId, pollOptionId, createdAt) VALUES ($1, $2, $3, $4)",
			sessionId, pollId, pollOption.PollOptionId, time.Now())

		if err != nil {
			http.Error(w, "Failed to save vote", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

        result, err := redisClient.ZIncrBy(ctx, pollId, 1.0, pollOption.PollOptionId).Result()

		if err != nil {
			http.Error(w, "Failed to Increment vote count in Redis", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		vote.Publish(pollId, util.Message{
			PollOptionID: pollOption.PollOptionId,
			Votes:        result, 		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}
