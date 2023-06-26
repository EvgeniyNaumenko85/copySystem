package repository

import (
	"copySys/db"
	"copySys/models"
	"copySys/pkg/logger"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FilePostgres struct {
	db *sqlx.DB
}

func NewFilePostgres() *FilePostgres {
	return &FilePostgres{}
}

func checkFileExist(fullPathToFile string) error {
	_, err := os.Stat(fullPathToFile)
	if os.IsNotExist(err) {
		return nil // файл не найден (true)
	} else {
		return models.ErrFileAlreadyExists
	}
}

func getFileIDByFileName(fileName string) (int, error) {

	var fileID int
	err := db.GetDBConn().QueryRow(db.GetFileIDByFileNameSql, fileName).Scan(&fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, models.ErrFileNotExists
		} else {
			logger.Error.Println(err.Error())
			return 0, err
		}
	}
	return fileID, nil
}

func getFileSizeLimitSql(userName string) (fileSizeLimit int, err error) {
	err = db.GetDBConn().QueryRow(db.CheckFileSizeLimitSql, userName).Scan(&fileSizeLimit)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	return fileSizeLimit, nil
}

func checkFileSizeToUpload(c *gin.Context) (handlerSizeInMb float64, err error) {
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}
	defer file.Close()

	var sizeOfOneMb int64 = 1024 * 1024
	handlerSizeInMb = float64(handler.Size / sizeOfOneMb)

	return handlerSizeInMb, nil
}

func getStorageFreeSpace(userName string) (freeSpace float64, err error) {
	err = db.GetDBConn().QueryRow(db.GetStorageFreeSpaceSql, userName).Scan(&freeSpace)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}
	return freeSpace, nil
}

func getFileSize(filePath string) (float64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	fileSize := float64(fileInfo.Size())

	var sizeOfOneMb float64 = 1024 * 1024
	fileSizeInMbString := fmt.Sprintf("%2f", math.Round(fileSize/sizeOfOneMb))
	fileSizeInMb, err := strconv.ParseFloat(fileSizeInMbString, 64)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	if fileSizeInMb == 0 {
		fileSizeInMb = 0.01
	}

	return fileSizeInMb, nil
}

func addFileInfoToDB(userId int, fileName, extension, path string, fileSize float64) (id int, err error) {
	err = db.GetDBConn().QueryRow(db.CreateFileSql, userId, fileName, extension, path, fileSize).Scan(&id)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func getFilePathByFileID(fileID int) (path string, err error) {
	err = db.GetDBConn().QueryRow(db.GetFilePathByFileIDSql, fileID).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", models.ErrFileNotExists
		} else {
			logger.Error.Println(err.Error())
			return "", err
		}
	}

	return path, nil
}

func getFileIDByFilePath(filePath string) (fileID int, err error) {
	err = db.GetDBConn().QueryRow(db.GetFileIDByFilePathSql, filePath).Scan(&fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, models.ErrFileNotExists
		} else {
			logger.Error.Println(err.Error())
			return 0, err
		}
	}

	return fileID, nil
}

func getAllFilesPaths() (filesPath []string, err error) {
	rows, err := db.GetDBConn().Query(db.GetAllFilesPathSql)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := models.File{}
		err = rows.Scan(
			&f.FilePath,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		filesPath = append(filesPath, f.FilePath)
	}

	if len(filesPath) == 0 {
		logger.Error.Println(models.ErrFilesNotExists)
		return filesPath, models.ErrFilesNotExists
	}

	return filesPath, nil
}

func getUserRoleByUserID(userID int) (userRole string, err error) {
	err = db.GetDBConn().QueryRow(db.GetUserRoleByIDSql, userID).Scan(&userRole)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", models.ErrFileNotExists
		} else {
			logger.Error.Println(err.Error())
			return "", err
		}
	}

	return userRole, err
}

func deleteFileInfoByFileID(fileID int) error {
	result, err := db.GetDBConn().Exec(db.DeleteFileByIDSql, fileID)
	foundRows, _ := result.RowsAffected()
	if foundRows == 0 {
		return models.ErrFileInfoNotFound
	}
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}

// outer func ===========>>
func (fp *FilePostgres) UploadFile(header *multipart.FileHeader, c *gin.Context) (int, error) {
	fileName := filepath.Base(header.Filename)
	partsOfFileName := strings.Split(fileName, ".")
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	fileName = partsOfFileName[0] + fileExtension
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	userName, err := getUserNameFromContext(c)
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	var folderToFileStore = "\\storage\\"

	fullPathToFile := currentDir + folderToFileStore + userName + "\\" + fileName

	err = checkFileExist(fullPathToFile)
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	userId, err := findUserIdByName(userName)
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	fileSizeLimit, err := getFileSizeLimitSql(userName)
	if err != nil {
		//todo log
		return 0, err
	}

	handlerSizeInMb, err := checkFileSizeToUpload(c)
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	freeSpace, err := getStorageFreeSpace(userName)
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	if handlerSizeInMb > float64(fileSizeLimit) {
		logger.Error.Println(err)
		return 0, models.ErrFileToBig
	}

	if freeSpace-handlerSizeInMb <= 0 {
		logger.Error.Println(err)
		return 0, models.ErrFileToBig
	}

	if err = c.SaveUploadedFile(header, fullPathToFile); err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	fileSize, err := getFileSize(fullPathToFile)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	localPathToFile := "." + "\\" + folderToFileStore + userName + "\\" + fileName

	fileId, err := addFileInfoToDB(userId, fileName, fileExtension, localPathToFile, fileSize)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	err = addAccessInfoToDB(fileId, userId)
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}
	return fileId, nil
}

func (fp *FilePostgres) GetFileByID(fileID int, userName string) (filePath string, err error) {

	userID, err := findUserIdByName(userName)
	if err != nil {
		logger.Error.Println(err)
		return "", err
	}

	err = checkUserToFileAccess(fileID, userID)
	if err != nil {
		logger.Error.Println(err)
		return "", err
	}

	filePath, err = getFilePathByFileID(fileID)
	if err != nil {
		logger.Error.Println(err)
		return "", err
	}

	return filePath, err

}

func (fp *FilePostgres) AllFilesInfo() (files []models.File, err error) {
	rows, err := db.GetDBConn().Query(db.AllFilesInfoSql)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := models.File{}
		err = rows.Scan(
			&file.ID,
			&file.UserId,
			&file.FileName,
			&file.Extension,
			&file.FileSize,
			&file.FilePath,
			&file.Added,
		)
		if err != nil {
			logger.Error.Println(err)
			continue
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		logger.Error.Println(models.ErrFilesNotExists)
		return files, models.ErrFilesNotExists
	}

	return files, nil
}

func (fp *FilePostgres) ShowAllUserFilesInfo(c *gin.Context) (files []models.File, err error) {
	userName, err := getUserNameFromContext(c)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	userID, err := findUserIdByName(userName)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	rows, err := db.GetDBConn().Query(db.GetAllUserFilesSql, userID)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := models.File{}
		err = rows.Scan(
			&file.ID,
			&file.UserId,
			&file.FileName,
			&file.Extension,
			&file.FileSize,
			&file.FilePath,
			&file.Added,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return files, models.ErrNoRows
	}
	return files, nil
}

func (fp *FilePostgres) FindFileByFileName(fileName, userName string) (file models.File, err error) {
	userID, err := findUserIdByName(userName)
	if err != nil {
		logger.Error.Println(err.Error())
		return file, err
	}

	fileID, err := getFileIDByFileName(fileName)
	if err != nil {
		logger.Error.Println(err.Error())
		return file, err
	}

	err = checkUserToFileAccess(fileID, userID)
	if err != nil {
		fmt.Println("err", err)
		logger.Error.Println(err.Error())
		return file, err
	}

	row := db.GetDBConn().QueryRow(db.GetFileByFileIDSql, fileID)
	err = row.Scan(
		&file.ID,
		&file.UserId,
		&file.FileName,
		&file.Extension,
		&file.FileSize,
		&file.FilePath,
		&file.Added,
	)
	if err != nil {
		logger.Error.Println(err.Error())
		return file, err
	}

	return file, err
}

func (fp *FilePostgres) DeleteFileByID(fileID int) error {
	filePath, err := getFilePathByFileID(fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return models.ErrNoRows
	}

	err = os.Remove(filePath)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	err = deleteFileInfoByFileID(fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	err = deleteAccessByFileID(fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}

func (fp *FilePostgres) DeleteAllFiles() (err error) {
	filesPath, err := getAllFilesPaths()
	if err != nil {
		logger.Error.Println(err.Error())
		return models.ErrNoRows
	}

	for i := 0; i < len(filesPath); i++ {
		fileID, err := getFileIDByFilePath(filesPath[i])
		if err != nil {
			logger.Error.Println(err.Error())
			return err
		}

		err = os.Remove(filesPath[i])
		if err != nil {
			logger.Error.Println(err.Error())
			return err
		}

		err = deleteFileInfoByFileID(fileID)
		if err != nil {
			logger.Error.Println(err.Error())
			return err
		}

		err = deleteAccessByFileID(fileID)
		if err != nil {
			logger.Error.Println(err.Error())
			return err
		}

	}

	return nil
}
