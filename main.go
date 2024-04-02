package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "go-api/controllers"
    "go-api/ws"

)


func main(){
    fmt.Println("Http server running")

    http.HandleFunc("/",controllers.Home_Controller)
    http.HandleFunc("/polls",controllers.Polls_Controller)
    http.HandleFunc("/polls/{pollId}", controllers.Poll_Controller)
    http.HandleFunc("/polls/{pollId}/votes",controllers.VotePoll_Controller)
    http.HandleFunc("/polls/{pollId}/results", ws.ResultsWebSocketHandler)


    fmt.Println("Server is listinig on port  http://localhost:3333")
    log.Fatal(http.ListenAndServe(":3333",nil))
}
