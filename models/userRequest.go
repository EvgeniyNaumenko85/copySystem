package models

type SingInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

//type RequestInput struct {
//	FileName string `json:"file_name,omitempty"`
//	UserName string `json:"user_name,omitempty"`
//}

type AccessRequest struct {
	FileID         int `json:"file_id"`
	AccessToUserID int `json:"user_id"`
}
