package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"html/template"

	"go-api/util"

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

    response,err := util.GetPoll(w,pollId)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PollHtlm_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not", http.StatusMethodNotAllowed)
		return
	}

	pollId := r.URL.Query().Get("pollId")
	if pollId == "" {
		http.Error(w, "Missing pollId parameter", http.StatusBadRequest)
		return
	}
    _, err := uuid.Parse(pollId)
    if err != nil {
        http.Error(w, "Invalid pollID parameter", http.StatusBadRequest)
        return
    }

    response,err := util.GetPoll(w,pollId)
	if err != nil {
		log.Fatal(err)
		return
	}


    tmpl, err := template.ParseFiles("templates/poll.html")
    if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		log.Fatal(err)
		return
    }

    w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, response)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
}
