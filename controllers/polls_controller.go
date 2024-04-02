package controllers

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "go-api/models"

     _ "github.com/lib/pq"
    "github.com/google/uuid"
)

func Polls_Controller (w http.ResponseWriter, r *http.Request){
        if r.Method != "POST" {
            http.Error(w, "Method not allowed",http.StatusMethodNotAllowed)
            return
        }

        var poll struct {
            Title string `json:"title"`
            Options []string `json:"options"`
        }

        decoder := json.NewDecoder(r.Body)
        if err := decoder.Decode(&poll); err != nil {
            http.Error(w,"Invalid requst body", http.StatusBadRequest)
        }

        if poll.Title == "" {
            http.Error(w, "Title cannot be empty", http.StatusBadRequest)
            return
        }

        if len(poll.Options) <2 {
            http.Error(w,"At least 2 options are required",http.StatusBadRequest)
            return
        }

        pollId := uuid.New().String()
        createdAt := time.Now()
        updatedAt := createdAt

        db, err := models.ConnectDB()
        if err != nil {
            http.Error(w,"Internal server error", http.StatusInternalServerError)
            log.Fatal(err)
            return
        }

        defer db.Close()

        tx, err := db.Begin()
        if err != nil {
            http.Error(w,"Faild to start transaction",http.StatusInternalServerError)
            log.Fatal(err)
            return
        }

        defer func(){
            if err != nil {
                tx.Rollback()
                return
            }
            err = tx.Commit()
        }()

        _, err = db.Exec("INSERT INTO poll (id, title, createdat, updatedat) VALUES ($1,$2,$3,$4)",
        pollId, poll.Title,createdAt,updatedAt)

        if err != nil {
            http.Error(w,"Faild to create poll", http.StatusInternalServerError)
            log.Fatal(err)
            return
        }

        for _, option := range poll.Options {
            optionId := uuid.New().String()
            _, err = tx.Exec("INSERT INTO polloption (id,title,pollid) VALUES ($1,$2,$3)",
                optionId,option,pollId)

            if err != nil {
                 http.Error(w,"Failed to insert option",http.StatusInternalServerError)
                 log.Fatal(err)
                 return
            }
        }

        createdPoll := models.Poll{
            Id: pollId,
            Title: poll.Title,
            CreatedAt: createdAt,
            UpdatedAt: updatedAt,
        }

        w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(createdPoll.Id)
    }
