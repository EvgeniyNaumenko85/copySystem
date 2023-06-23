package models

type UserRequest struct {
	OldUserId int `json:"old_user_id"`
	NewUserId int `json:"new_user_id"`
}

type SingInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RequestInput struct {
	FileName string `json:"file_name,omitempty"`
	UserName string `json:"user_name,omitempty"`
}
