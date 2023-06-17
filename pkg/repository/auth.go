package repository

import (
	"tasks_API/db"
	"tasks_API/models"
	"tasks_API/pkg/logger"
)

type AuthPostgres struct {
	//repo *repository.Repository
}

func NewAuthPostgres() *AuthPostgres {
	return &AuthPostgres{}
}

func (ap *AuthPostgres) GetUser(username, password, role string) (u models.User, err error) {
	row := db.GetDBConn().
		QueryRow("select id, username, email, role from users WHERE username = $1 AND password_hash = $2",
			username, password)

	err = row.Scan(&u.Id, &u.UserName, &u.Email, &u.Role)
	if err != nil {
		logger.Error.Println(err.Error())
		return models.User{}, err
	}

	return u, nil

}

func (ap *AuthPostgres) CreateUser(u models.User) (id int, err error) {
	err = db.GetDBConn().
		QueryRow("insert into users (username, email, password_hash, role) values ($1, $2, $3, $4) RETURNING id",
			u.UserName, u.Email, u.Password, u.Role).Scan(&id)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	return id, nil
}
