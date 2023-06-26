package db

const (
	GetAllTasksSql                = "SELECT id, user_id, task_name, is_done, description, added, deadLine FROM tasks"
	GetTaskByIDSql                = "SELECT id, user_id, task_name, is_done, description, added, deadLine FROM tasks WHERE id = $1"
	GetTaskByUserIDSqlFunc        = "SELECT tasks.id, user_id, task_name, is_done, description, added, deadLine FROM tasks JOIN users ON tasks.user_id = users.id AND tasks.user_id = $1"
	GetUndoneTasksByUserIDSqlFunc = "SELECT * FROM getUndoneTasksByUserId($1)"
	GetOverdueTasksSqlFunc        = "SELECT * FROM get_overdue_tasks($1)"
	CreateTaskSql                 = "INSERT INTO tasks (user_id, task_name, is_done, description, deadline) VALUES ($1, $2, $3, $4, $5)  RETURNING id"
	UpdateTaskByIDSql             = "UPDATE tasks SET task_name = $1, description = $2, is_done = $3, deadline=$4 WHERE id = $5"
	DeleteTaskByIDSql             = "DELETE FROM tasks WHERE id = $1"
	ReassignTaskSqlProcedure      = "CALL reassign_task($1, $2, $3)"
)

// users
const (
	GetAllUsersSql            = "SELECT id, user_name, email, user_role FROM users"
	GetAllUsersIDsSql         = "SELECT id FROM users"
	GetUserByIDSql            = "SELECT id, user_name, email, user_role  FROM users WHERE id = $1"
	GetUserRoleByIDSql        = "SELECT user_role  FROM users WHERE id = $1"
	UpdateUserByIDSql         = "UPDATE users SET id=$1, user_name = $2, email = $3, user_role = $4 WHERE id = $5"
	DeleteUserByIDSql         = "DELETE FROM users WHERE id = $1"
	GetIdUserByNameSql        = "SELECT id FROM users WHERE user_name = $1"
	CheckUserExistByUserIDSql = "SELECT id FROM users WHERE id = $1"
)

// access
const (
	CreateAccessSql         = "INSERT INTO access (user_id, file_id) VALUES ($1, $2)"
	CheckAccessInTableSql   = "SELECT id FROM access WHERE (file_id = $1 AND user_id =$2)"
	DeleteAccessByFileIDSql = "DELETE FROM access WHERE id = $1"
	DeleteUserAccessSql     = "DELETE FROM access WHERE (file_id = $1 AND user_id =$2 AND user_id !=$3)"
	DeleteAccessToAllSql    = "DELETE FROM access WHERE (file_id = $1 AND user_id !=$2)"
)

// files
const (
	GetAllUserFilesSql     = "SELECT id, user_id, file_name, extension, file_path, file_size, added FROM files WHERE user_id = $1"
	GetFileByFileIDSql     = "SELECT id, user_id, file_name, extension, file_path, file_size, added FROM files WHERE id = $1"
	AllFilesInfo           = "SELECT id, user_id, file_name, extension, file_size, file_path, added FROM files"
	CreateFileSql          = "INSERT INTO files (user_id, file_name, extension, file_path, file_size) VALUES ($1, $2, $3, $4, $5)  RETURNING id"
	GetFilePathByFileIDSql = "SELECT file_path FROM files WHERE id = $1 "
	GetFileIDByFilePathSql = "SELECT id FROM files WHERE file_path = $1 "
	GetFileIDByFileNameSql = "SELECT id FROM files WHERE file_name = $1"
	GetAllFilesPath        = "SELECT file_path FROM files"
	CheckFileSizeLimitSql  = "SELECT file_size_lim FROM users WHERE user_name =$1"
	CheckFileByFileIDSql   = "SELECT id FROM files WHERE id = $1"
	DeleteFileByIDSql      = "DELETE FROM files WHERE id = $1"
	GetStorageFreeSpaceSql = "SELECT  users.storage_size_lim - COALESCE(SUM(files.file_size), 0) AS remaining_storage FROM users LEFT JOIN files ON files.user_id = users.id WHERE users.user_name = $1 GROUP BY users.id, users.storage_size_lim"
)

const (
	SetLimitsToUser = "UPDATE users SET file_size_lim = $1, storage_size_lim = $2 WHERE id = $3"
)

//
//SELECT
//users.id AS user_id,
//users.storage_size_lim - COALESCE(SUM(files.file_size), 0) AS remaining_storage
//FROM
//users
//LEFT JOIN
//files ON files.user_id = users.id
//WHERE
//users.id = <ваш user_id>
//GROUP BY
//users.id, users.storage_size_lim;
