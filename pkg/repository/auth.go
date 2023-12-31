package repository

import (
	"copySys/db"
	"copySys/models"
	"copySys/pkg/logger"
)

type AuthPostgres struct {
	//repo *repository.Repository
}

func NewAuthPostgres() *AuthPostgres {
	return &AuthPostgres{}
}

func (ap *AuthPostgres) GetUser(username, password, role string) (u models.User, err error) {
	row := db.GetDBConn().
		QueryRow("select id, user_name, email, user_role from users WHERE user_name = $1 AND password_hash = $2",
			username, password)

	err = row.Scan(&u.ID, &u.UserName, &u.Email, &u.Role)
	if err != nil {
		logger.Error.Println("GetUser func: ", err.Error())
		return models.User{}, err
	}

	return u, nil

}

func (ap *AuthPostgres) CreateUser(u models.User) (id int, err error) {
	err = db.GetDBConn().
		QueryRow("insert into users (user_name, email, password_hash, user_role) values ($1, $2, $3, $4) RETURNING id",
			u.UserName, u.Email, u.Password, u.Role).Scan(&id)
	if err != nil {
		logger.Error.Println("CreateUser func: ", err.Error())
		return 0, err
	}

	return id, nil
}
