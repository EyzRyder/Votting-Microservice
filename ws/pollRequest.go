package ws

import (
    "log"
    "net/http"

   // "go-api/util"

    "github.com/gorilla/websocket"
	"github.com/google/uuid"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func ResultsWebSocketHandler(w http.ResponseWriter, r *http.Request) {
upgrader.CheckOrigin = func(r *http.Request) bool {
    // Check if the request origin is allowed
    allowedOrigins := []string{"https://hoppscotch.io", "http://localhost:3000"}
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

    if pollID == "" {
        log.Println("Missing pollId parameter")
        return
    }



}
