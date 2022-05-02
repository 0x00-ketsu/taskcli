package model

import "time"

type Task struct {
	ID          int64     `storm:"id,increment" json:"id"`
	Title       string    `storm:"index" json:"title"`
	Content     string    `storm:"index" json:"content"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"create_at"`
	IsCompleted bool      `json:"is_completed"`
	IsDeleted   bool      `json:"is_deleted"`
}
