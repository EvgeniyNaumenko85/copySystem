package db

//const CreateDbFiles_db = "CREATE DATABASE files_db;"

const (
	CreateTableUsers = `CREATE TABLE IF NOT EXISTS users
(
id             	BIGSERIAL PRIMARY KEY,
user_name       VARCHAR(30) NOT NULL,
email         	VARCHAR(30) NOT NULL UNIQUE,
password_hash 	VARCHAR(255) NOT NULL,
role  			VARCHAR(30) NOT NULL
);`
	CreateTableFiles = `CREATE TABLE IF NOT EXISTS files
(
id               BIGSERIAL PRIMARY KEY,
user_id          integer REFERENCES users (id) ON DELETE CASCADE,
file_name        VARCHAR(30)  NOT NULL UNIQUE,
extension        VARCHAR(10),
path             VARCHAR(255) NOT NULL UNIQUE,
description      VARCHAR(255),
deleted          BOOLEAN      NOT NULL DEFAULT false,
added            TIMESTAMP    NOT NULL DEFAULT now(),
version          INTEGER NOT NULL --DEFAULT  --to do: при отправке одноименного файла (с уже существующим file_name) увеличивалась версия файла
);`

	CreateTableShedules = `CREATE TABLE IF NOT EXISTS shedules
(
    id               BIGSERIAL PRIMARY KEY,
    user_id          INTEGER REFERENCES users (id) ON DELETE CASCADE,
    file_id          INTEGER REFERENCES files (id) ON DELETE CASCADE,
    copy_time        TIMESTAMP
);`
)

const (
	DropShedulesTable = `DROP TABLE IF EXISTS shedules;`
	DropFilesTable    = `DROP TABLE IF EXISTS files;`
	DropUsersTable    = `DROP TABLE IF EXISTS users;`
)
