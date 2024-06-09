package main

import (
	"fmt"
	"log"
	"net/http"

	"go-api/controllers"
	"go-api/util"
	"go-api/ws"
)

func main() {
	fmt.Println("Http server running")

	vote := util.NewVotingPubSub()

	http.HandleFunc("/", controllers.Home_Controller)
	http.HandleFunc("/createPoll", controllers.PollForm_Controller)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/polls", controllers.Polls_Controller)
	http.HandleFunc("/polls/{pollId}", controllers.Poll_Controller)
	http.HandleFunc("/pollsHtml/{pollId}", controllers.PollHtlm_Controller)
	http.HandleFunc("/polls/{pollId}/votes", controllers.VotePoll_Controller(vote))
	http.HandleFunc("/polls/{pollId}/results", ws.ResultsWebSocketHandler(vote))

	fmt.Println("Server is listinig on port  http://localhost:3333")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
