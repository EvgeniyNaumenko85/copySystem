package models

import "time"

type File struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	FileName  string    `json:"file_name"`
	Extension string    `json:"extension"`
	FileSize  string    `json:"file_size"`
	Deleted   bool      `json:"-" `
	Added     time.Time `json:"added"`
}
