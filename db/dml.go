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

const (
	GetAllUsersSql    = "SELECT id, username, email, role FROM users"
	GetUserByIDSql    = "SELECT id, username, email, role  FROM users WHERE id = $1"
	UpdateUserByIDSql = "UPDATE users SET username = $1, email = $2, role = $3 WHERE id = $4"
	DeleteUserByIDSql = "DELETE FROM users WHERE id = $1"

	//ChekUserInTableSql = "SELECT EXISTS id FROM users WHERE id = $1"
	//AddUserSql        = "INSERT INTO users (id, name, email, date_of_Birth, surname) VALUES ($1, $2, $3, $4, $5) RETURNING id"
)
