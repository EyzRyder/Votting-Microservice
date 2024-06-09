package models

import (
	"time"
)

type Poll struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PollOption struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	PollId string `json:"pollId"`
}

type Vote struct {
	Id           string    `json:"id"`
	SessionId    string    `json:"sessionId"`
	PollId       string    `json:"pollId"`
	PollOptionId string    `json:"pollOptionId"`
	CreatedAt    time.Time `json:"createdAt"`
}
