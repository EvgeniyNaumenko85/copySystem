package models

import "errors"

var ErrNoRows = errors.New("sql: no rows in result set")
var ErrNoRowsSQL = "sql: no rows in result set"

var ErrUserNotFound = errors.New("sql: user not found")
var ErrInvalidEmailForm = errors.New("invalid email form")
var ErrNotUnicUser = "pq: повторяющееся значение ключа нарушает ограничение уникальности \"users_email_key\""
var ErrNotUnicUserName = "pq: повторяющееся значение ключа нарушает ограничение уникальности \"users_user_name_key\""
var ErrUserNotExists = errors.New("user is not exists in DB")
var ErrCantGetUserName = errors.New("can't  get user name")

var ErrFileAlreadyExists = errors.New("file already exists")
var ErrFileAccessDenied = errors.New("file access denied")
var ErrFileToBig = errors.New("file to upload is too big")
var ErrFileNotExists = errors.New("file is not exists in DB")
var ErrFileInfoNotFound = errors.New("file info is not found")

var ErrAccessInfoNotFound = errors.New("access info is not found")

var ErrCantGetInfoFromHeader = errors.New("can't  get info from Headers")
