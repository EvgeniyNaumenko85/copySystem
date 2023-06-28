package db

// users
const (
	GetAllUsersSql            = "SELECT id, user_name, email, user_role FROM users"
	GetUserByIDSql            = "SELECT id, user_name, email, user_role  FROM users WHERE id = $1"
	GetUserRoleByIDSql        = "SELECT user_role  FROM users WHERE id = $1"
	GetUserNameByUserID       = "SELECT user_name  FROM users WHERE id = $1"
	UpdateUserByIDSql         = "UPDATE users SET id=$1, user_name = $2, email = $3, user_role = $4 WHERE id = $5"
	DeleteUserByIDSql         = "DELETE FROM users WHERE id = $1"
	GetUserIDByUserNameSql    = "SELECT id FROM users WHERE user_name = $1"
	CheckUserExistByUserIDSql = "SELECT id FROM users WHERE id = $1"
)

// access
const (
	CreateAccessSql         = "INSERT INTO access (user_id, file_id) VALUES ($1, $2)"
	CheckAccessInTableSql   = "SELECT id FROM access WHERE (file_id = $1 AND user_id =$2)"
	DeleteAccessByFileIDSql = "DELETE FROM access WHERE id = $1"
	DeleteUserAccessSql     = "DELETE FROM access WHERE (file_id = $1 AND user_id =$2 AND user_id !=$3)"
	DeleteAccessToAllSql    = "DELETE FROM access WHERE (file_id = $1 AND user_id !=$2)"
	RemoveAccessToAllSql    = "UPDATE files SET access_to_all = false  WHERE id = $1 AND access_to_all = true"
)

// files
const (
	GetAllUserFilesSql              = "SELECT id, user_id, file_name, extension, file_path, file_size, added FROM files WHERE user_id = $1"
	GetFileByFileIDSql              = "SELECT id, user_id, file_name, extension, file_path, file_size, added FROM files WHERE id = $1"
	AllFilesInfoSql                 = "SELECT id, user_id, file_name, extension, file_size, file_path, added FROM files"
	CreateFileSql                   = "INSERT INTO files (user_id, file_name, extension, file_path, file_size) VALUES ($1, $2, $3, $4, $5)  RETURNING id"
	GetFilePathByFileIDSql          = "SELECT file_path FROM files WHERE id = $1 "
	GetFileIDByFilePathSql          = "SELECT id FROM files WHERE file_path = $1 "
	GetFileIDByFileNameSql          = "SELECT id FROM files WHERE file_name = $1"
	GetAllFilesPathSql              = "SELECT file_path FROM files"
	CheckFileSizeLimitSql           = "SELECT file_size_lim FROM users WHERE user_name =$1"
	CheckFileByFileIDSql            = "SELECT id FROM files WHERE id = $1"
	DeleteFileByIDSql               = "DELETE FROM files WHERE id = $1"
	GetFilesQuantityByUserIDSql     = "SELECT COUNT(*) AS row_count FROM files WHERE user_id = $1;"
	GetStorageFreeSpaceSql          = "SELECT  users.storage_size_lim - COALESCE(SUM(files.file_size), 0) AS remaining_storage FROM users LEFT JOIN files ON files.user_id = users.id WHERE users.user_name = $1 GROUP BY users.id, users.storage_size_lim"
	SetFileAccessToAllUsers         = "UPDATE files SET access_to_all = true  WHERE id = $1 AND access_to_all = false"
	CheckAccessToAllInTableFilesSql = "SELECT id FROM files WHERE (id = $1 AND  access_to_all = true)"
)

// limits
const (
	SetLimitsToUserSql = "UPDATE users SET file_size_lim = $1, storage_size_lim = $2 WHERE id = $3"
)
