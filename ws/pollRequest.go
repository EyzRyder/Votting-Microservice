package ws

import (
	"log"
	"net/http"

	"go-api/util"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ResultsWebSocketHandler(vote *util.VotingPubSub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Keep-Alive", "timeout=600")

		upgrader.CheckOrigin = func(r *http.Request) bool {
			allowedOrigins := []string{"https://hoppscotch.io"}
			origin := r.Header.Get("Origin")
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == origin {
					return true
				}
			}
			return false
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		defer conn.Close()

		pollID := r.PathValue("pollId")

		if pollID == "" {
			http.Error(w, "Missing pollId parameter", http.StatusBadRequest)
			return
		}
		_, err = uuid.Parse(pollID)
		if err != nil {
			http.Error(w, "Invalid pollID parameter", http.StatusBadRequest)
			return
		}

		vote.Subscribe(pollID, func(msg util.Message) {
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("Error writing message to WebSocket: ", err)
				return
			}

		})
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading from WebSocket:", err)
				break
			}
		}

	}
}
