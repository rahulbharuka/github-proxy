package model

import "time"

// Comment is a model for comment
type Comment struct {
	Author    string    `json:"author"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
