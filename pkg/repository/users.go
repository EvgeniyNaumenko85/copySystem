package repository

import (
	"github.com/jmoiron/sqlx"
	"tasks_API/db"
	"tasks_API/models"
	"tasks_API/pkg/logger"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres() *UserPostgres {
	return &UserPostgres{}
}

func (up *UserPostgres) GetAllUsers() (users []models.User, err error) {
	rows, err := db.GetDBConn().Query(db.GetAllUsersSql)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := models.User{}
		err = rows.Scan(
			&u.Id,
			&u.UserName,
			&u.Email,
			&u.Role,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

func (up *UserPostgres) GetUserByID(id int) (user models.User, err error) {
	var u models.User
	row := db.GetDBConn().
		QueryRow(db.GetUserByIDSql, id)
	err = row.Scan(
		&u.Id,
		&u.UserName,
		&u.Email,
		&u.Role,
	)
	if err != nil {
		logger.Error.Println(err.Error())
		return u, err
	}

	return u, nil
}

func (up *UserPostgres) UpdateUserByID(id int, u models.User) (err error) {
	result, err := db.GetDBConn().Exec(db.UpdateUserByIDSql, u.UserName, u.Email, u.Role, id)

	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	updatedRows, _ := result.RowsAffected()
	if updatedRows == 0 {

		return models.ErrUserNotFound
	}

	return nil
}

func (up *UserPostgres) DeleteUserByID(id int) (err error) {
	result, err := db.GetDBConn().Exec(db.DeleteUserByIDSql, id)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	deletedRows, _ := result.RowsAffected()
	if deletedRows == 0 {
		return models.ErrNoRows
	}

	return nil
}
