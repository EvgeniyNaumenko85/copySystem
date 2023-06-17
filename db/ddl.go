package db

//const CreateDbTasks_db = "CREATE DATABASE tasks_db;"

const (
	CreateTableUsers = `CREATE TABLE IF NOT EXISTS users
(
id             BIGSERIAL PRIMARY KEY,
-- full_name      VARCHAR(30) NOT NULL,
username       VARCHAR(30) NOT NULL,
email         VARCHAR(30) NOT NULL UNIQUE,
password_hash VARCHAR(255) NOT NULL,
role  VARCHAR(30) NOT NULL
-- date_of_Birth DATE NOT NULL
);`
	CreateTableTasks = `CREATE TABLE IF NOT EXISTS tasks
(
id               BIGSERIAL PRIMARY KEY,
user_id          integer REFERENCES users (id) ON DELETE CASCADE,
task_name        VARCHAR(30)  NOT NULL,
description      VARCHAR(255) NOT NULL,
is_done          BOOLEAN      NOT NULL,
added            TIMESTAMP    NOT NULL DEFAULT now(),
deadline         TIMESTAMP,
done_at          TIMESTAMP
);`
)

const (
	DropTasksTable = `DROP TABLE IF EXISTS tasks;`
	DropUsersTable = `DROP TABLE IF EXISTS users;`
)
