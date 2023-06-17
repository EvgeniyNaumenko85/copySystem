package models

import "time"

// Task defines the structure for an API tasks
type Task struct {
	ID          int       `json:"id"`
	UserId      int       `json:"user_id,omitempty"`
	Name        string    `json:"name"`
	Done        bool      `json:"done"`
	Description string    `json:"description"`
	Added       time.Time `json:"added"`
	DeadLine    time.Time `json:"deadline,omitempty"`
	DoneAt      time.Time `json:"done_at,omitempty"`
}

// Tasks is a collection of Task
type Tasks []*Task
