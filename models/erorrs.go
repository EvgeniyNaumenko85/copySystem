package models

import "errors"

var ErrNoRows = errors.New("sql: no rows in result set")
var ErrUserNotFound = errors.New("sql: user not found")
var ErrNotUnicUser = "pq: повторяющееся значение ключа нарушает ограничение уникальности \"users_email_key\""
var ErrNoRowsSQL = "sql: no rows in result set"

var ErrFileAlredyExists = errors.New("file already exists")
