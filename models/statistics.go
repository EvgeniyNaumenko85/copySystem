package models

type UserStat struct {
	UserId        int     `json:"user_id"`
	UserName      string  `json:"user_name"`
	FilesQuantity int     `json:"files_qty"`
	FreeSpace     float64 `json:"free_space"`
}
