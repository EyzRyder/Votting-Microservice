package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "go-api/controllers"
    "go-api/ws"

    "github.com/gorilla/websocket"

)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
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
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading to WebSocket:", err)
        return
    }
    defer conn.Close()

    // Send "Hello, world!" message at 30-second intervals
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Send message to WebSocket connection
            if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, world!")); err != nil {
                log.Println("Error writing message to WebSocket:", err)
                return
            }
        }
    }
}

func main(){
    fmt.Println("Http server running")

    http.HandleFunc("/",controllers.Home_Controller)
    http.HandleFunc("/hello",helloWorldHandler)
    http.HandleFunc("/polls",controllers.Polls_Controller)
    http.HandleFunc("/polls/{pollId}", controllers.Poll_Controller)
    http.HandleFunc("/polls/{pollId}/votes",controllers.VotePoll_Controller)
    http.HandleFunc("/polls/{pollId}/results", ws.ResultsWebSocketHandler)


    fmt.Println("Server is listinig on port  http://localhost:3333")
    log.Fatal(http.ListenAndServe(":3333",nil))
}
