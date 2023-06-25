package repository

import (
	"copySys/db"
	"copySys/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type LimitsPostgres struct {
	db *sqlx.DB
}

func NewLimitsPostgres() *LimitsPostgres {
	return &LimitsPostgres{}
}

func (lp *LimitsPostgres) SetLimitsToUser(userID, fileSizeLim int) (err error) {
	err = checkUserExistByUserID(userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	_, err = db.GetDBConn().Exec(db.SetLimitsToUser, fileSizeLim, userID)
	if err != nil {
		logger.Error.Println(err.Error())
		fmt.Println(err)
	}

	return nil
}
