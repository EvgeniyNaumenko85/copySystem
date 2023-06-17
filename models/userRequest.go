package models

type UserRequest struct {
	OldUserId int `json:"old_user_id"`
	NewUserId int `json:"new_user_id"`
}
