package repository

import (
	"copySys/db"
	"copySys/models"
	"copySys/pkg/logger"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres() *UserPostgres {
	return &UserPostgres{}
}

func getUserNameFromContext(c *gin.Context) (string, error) {
	userNameTypeAny, ok := c.Get("userName")
	if !ok {
		logger.Error.Println(models.ErrCantGetUserName.Error())
		return "", models.ErrCantGetUserName
	} else {
		userName := fmt.Sprintf("%v", userNameTypeAny)
		return userName, nil
	}
}

func findUserIdByName(userName string) (int, error) {
	if userName == "" {
		return 0, models.ErrUserNotExists
	}
	var ID int
	err := db.GetDBConn().QueryRow(db.GetIdUserByNameSql, userName).Scan(&ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, models.ErrUserNotExists
		} else {
			logger.Error.Println(err.Error())
			return 0, err
		}
	}
	return ID, nil
}

// GetAllUsers outer func ===========>>
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
			&u.ID,
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
		&u.ID,
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
	result, err := db.GetDBConn().Exec(db.UpdateUserByIDSql, u.ID, u.UserName, u.Email, u.Role, id)

	if err != nil {
		logger.Error.Println("UpdateUserByID func: ", err.Error())
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
