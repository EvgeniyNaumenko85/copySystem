package repository

import (
	"copySys/db"
	"copySys/models"
	"copySys/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type StatisticsPostgres struct {
	db *sqlx.DB
}

func NewStatisticsPostgres() *StatisticsPostgres {
	return &StatisticsPostgres{}
}

func (sp *StatisticsPostgres) GetUserStatistics(userID int) (userStat models.UserStat, err error) {
	err = checkUserExistByUserID(userID)
	if err != nil {
		logger.Error.Println(err)
		return userStat, err
	}

	var userName string
	err = db.GetDBConn().QueryRow(db.GetUserNameByUserID, userID).Scan(&userName)
	if err != nil {
		logger.Error.Println(err)
		return userStat, err
	}

	var filesQuantity int
	err = db.GetDBConn().QueryRow(db.GetFilesQuantityByUserIDSql, userID).Scan(&filesQuantity)
	if err != nil {
		logger.Error.Println(err)
		fmt.Println(err)
	}

	freeSpace, err := getStorageFreeSpace(userName)
	if err != nil {
		logger.Error.Println(err)
		return userStat, err
	}

	userStat.UserId = userID
	userStat.UserName = userName
	userStat.FilesQuantity = filesQuantity
	userStat.FreeSpace = freeSpace

	return userStat, nil
}
