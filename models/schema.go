package models

import (
    "time"
)

type Poll struct {
    Id string `json:"id"`
    Title string `json:"title"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

