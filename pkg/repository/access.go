package repository

import (
	"copySys/db"
	"copySys/models"
	"copySys/pkg/logger"
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
		logger.Error.Println(err)
		return err
	}

	userRole, err := getUserRoleByUserID(userID)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	resultFromFilesTable, err := db.GetDBConn().Exec(db.CheckAccessToAllInTableFilesSql, fileID)
	foundRows, _ := resultFromFilesTable.RowsAffected()
	if foundRows == 0 && userRole != "admin" {
		result, err := db.GetDBConn().Exec(db.CheckAccessInTableSql, fileID, userID)
		if err != nil {
			logger.Error.Println(err)
			return err
		}
		foundRows, _ = result.RowsAffected()
		if foundRows == 0 {
			return models.ErrFileAccessDenied
		}
	}

	return nil
}

func addAccessInfoToDB(fileId, userId int) error {
	_, err := db.GetDBConn().Exec(db.CreateAccessSql, userId, fileId)
	if err != nil {
		logger.Error.Println(err)
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
		logger.Error.Println(err)
		return err
	}
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

func removeAccessToAllinAccessTable(fileId, userId int) error {
	_, err := db.GetDBConn().Exec(db.DeleteAccessToAllSql, fileId, userId)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

func removeAccessToAllInFilesTable(fileId int) error {
	_, err := db.GetDBConn().Exec(db.RemoveAccessToAllSql, fileId)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

func deleteAccessByFileID(fileID int) error {
	_, err := db.GetDBConn().Exec(db.DeleteAccessByFileIDSql, fileID)
	if err != nil {
		logger.Error.Println(err)
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
		logger.Error.Println(err)
		return err
	}

	return nil
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
		logger.Error.Println(err)
		return models.ErrFileAccessDenied
	}

	err = checkUserExistByUserID(userID)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	err = addAccessInfoToDB(fileID, accessToUserID)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

func setFileAccessToAllUsers(fileID int) error {
	result, err := db.GetDBConn().Exec(db.SetFileAccessToAllUsers, fileID)
	foundRows, _ := result.RowsAffected()
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if foundRows == 0 {
		return models.ErrFileAccessToAllUsers
	}

	return nil
}

func (ap *AccessPostgres) ProvidingAccessAll(userID, fileID int) (err error) {

	err = checkFileByFileID(fileID)
	if err != nil {
		logger.Error.Println(err)
		return models.ErrFileNotExists
	}

	err = checkUserToFileAccess(fileID, userID)
	if err != nil {
		logger.Error.Println(err)
		return models.ErrFileAccessDenied
	}

	err = setFileAccessToAllUsers(fileID)
	if err != nil {
		logger.Error.Println(err)
		return models.ErrFileAccessToAllUsers
	}

	return nil
}

func (ap *AccessPostgres) RemoveAccess(fileID, accessToUserID, userID int) error {
	err := checkUserToFileAccess(fileID, userID)
	if err != nil {
		logger.Error.Println(err)
		return models.ErrFileAccessDenied
	}

	err = checkUserExistByUserID(userID)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	err = removeAccessInfoFromDB(fileID, accessToUserID, userID)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

func (ap *AccessPostgres) RemoveAccessToAll(fileID, userID int) error {

	err := checkUserToFileAccess(fileID, userID)
	if err != nil {
		logger.Error.Println(err)
		return models.ErrFileAccessDenied
	}

	err = checkUserExistByUserID(userID)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	err = removeAccessToAllinAccessTable(fileID, userID)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	err = removeAccessToAllInFilesTable(fileID)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
