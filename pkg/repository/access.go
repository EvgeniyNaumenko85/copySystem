package repository

import (
	"copySys/db"
	"copySys/models"
	"copySys/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AccessPostgres struct {
	db *sqlx.DB
}

func NewAccessPostgres() *AccessPostgres {
	return &AccessPostgres{}
}

func checkUserToFileAccess(fileID, userID int) error {
	_, err := getFilePathByFileID(fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	userRole, err := getUserRoleByUserID(userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	result, err := db.GetDBConn().Exec(db.CheckAccessInTableSql, fileID, userID)
	foundRows, _ := result.RowsAffected()
	if foundRows == 0 && userRole != "admin" {
		return models.ErrFileAccessDenied
	}
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}

func addAccessInfoToDB(fileId, userId int) error {
	_, err := db.GetDBConn().Exec(db.CreateAccessSql, userId, fileId)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}

func removeAccessInfoFromDB(fileId, accessToUserID, userId int) error {
	result, err := db.GetDBConn().Exec(db.DeleteUserAccessSql, fileId, accessToUserID, userId)
	foundRows, _ := result.RowsAffected()
	if foundRows == 0 {
		return models.ErrFileAccessDenied
	}
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}

func deleteAccessByFileID(fileID int) error {
	_, err := db.GetDBConn().Exec(db.DeleteAccessByFileIDSql, fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}

func checkUserExistByUserID(userID int) error {
	result, err := db.GetDBConn().Exec(db.CheckUserExistByUserIDSql, userID)
	foundRows, _ := result.RowsAffected()
	if foundRows == 0 {
		return models.ErrFileAccessDenied
	}
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}

func getAllUsersIDs() (allUsersIDs []int, err error) {
	rows, err := db.GetDBConn().Query(db.GetAllUsersIDsSql)
	if err != nil {
		logger.Error.Println(err.Error())
		return allUsersIDs, err
	}
	defer rows.Close()

	for rows.Next() {
		u := models.User{}
		err = rows.Scan(
			&u.ID,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		allUsersIDs = append(allUsersIDs, u.ID)
	}

	return allUsersIDs, nil
}

func checkFileByFileID(fileID int) error {
	result, err := db.GetDBConn().Exec(db.CheckFileByFileIDSql, fileID)
	foundRows, _ := result.RowsAffected()
	if foundRows == 0 {
		logger.Error.Println(err)
		return models.ErrFileNotExists
	}

	return nil
}

// outer func ===========>>
func (ap *AccessPostgres) ProvidingAccess(fileID, accessToUserID, userID int) error {
	err := checkUserToFileAccess(fileID, userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return models.ErrFileAccessDenied
	}

	err = checkUserExistByUserID(userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	err = addAccessInfoToDB(fileID, accessToUserID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}

func (ap *AccessPostgres) ProvidingAccessAll(userID, fileID int) (err error) {

	err = checkFileByFileID(fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return models.ErrFileNotExists
	}

	err = checkUserToFileAccess(fileID, userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return models.ErrFileAccessDenied
	}

	allUsersIDs, err := getAllUsersIDs()
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	for i := 0; i < len(allUsersIDs); i++ {
		currentUserID := allUsersIDs[i]
		if currentUserID != userID {
			result, _ := db.GetDBConn().Exec(db.CheckAccessInTableSql, fileID, currentUserID)
			foundRows, _ := result.RowsAffected()
			if foundRows != 0 {
				continue
			} else {
				err = addAccessInfoToDB(fileID, currentUserID)
				if err != nil {
					logger.Error.Println(err.Error())
					fmt.Println(err)
					continue
				}
			}
		}
	}

	return nil
}

func (ap *AccessPostgres) RemoveAccess(fileID, accessToUserID, userID int) error {
	err := checkUserToFileAccess(fileID, userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return models.ErrFileAccessDenied
	}

	err = checkUserExistByUserID(userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	err = removeAccessInfoFromDB(fileID, accessToUserID, userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}
