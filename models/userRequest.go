package models

type SingInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AccessRequest struct {
	FileID         int `json:"file_id"`
	AccessToUserID int `json:"user_id"`
}

type LimitRequest struct {
	FileSizeLim int `json:"file_size_lim"`
}
