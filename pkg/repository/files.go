package repository

import (
	"copySys/db"
	"copySys/models"
	"copySys/pkg/logger"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
		//fmt.Println("err: ", err)
		return models.ErrFileAlreadyExists
	}
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

func checkFileSizeToUpload(fileSizeLimit int, c *gin.Context) error {
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file:", err)
		c.String(http.StatusBadRequest, "Error retrieving file")
		return err
	}
	defer file.Close()

	// Проверка размера файла
	var sizeOfOneMb int64 = 1024 * 1024
	handlerSizeInMb := handler.Size / sizeOfOneMb
	if handlerSizeInMb > int64(fileSizeLimit) {
		fmt.Println("handler.Size: ", handler.Size)
		//todo log
		return models.ErrFileToBig
	}

	return nil
}

func getFileSize(filePath string) (int, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	fileSize := float64(fileInfo.Size())

	var sizeOfOneMb float64 = 1024 * 1024
	fileSizeInMB := int(math.Ceil(fileSize / sizeOfOneMb))

	return fileSizeInMB, nil

}

func addFileInfoToDB(userId int, fileName, extension, path string, fileSize int) (id int, err error) {
	err = db.GetDBConn().QueryRow(db.CreateFileSql, userId, fileName, extension, path, fileSize).Scan(&id)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func addAccessInfoToDB(fileId, userId int) error {
	_, err := db.GetDBConn().Exec(db.CreateAccessSql, userId, fileId)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
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

func checkUserToFileAccess(fileID, userID int) error {
	_, err := getFilePathByFileID(fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	result, err := db.GetDBConn().Exec(db.CheckAccessInTableSql, fileID, userID)
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

//func deleteAccessByFileID(fileID int) error {
//	result, err := db.GetDBConn().Exec(db.DeleteAccessByFileIDSql, fileID)
//	foundRows, _ := result.RowsAffected()
//	if foundRows == 0 {
//		return models.ErrAccessInfoNotFound
//	}
//	if err != nil {
//		logger.Error.Println(err.Error())
//		return err
//	}
//
//	return nil
//}

func deleteAccessByFileID(fileID int) error {
	_, err := db.GetDBConn().Exec(db.DeleteAccessByFileIDSql, fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}

// outer func ===========>>

func (fp *FilePostgres) UploadFile(header *multipart.FileHeader, c *gin.Context) (int, error) {
	fileName := filepath.Base(header.Filename)
	fileExtension := filepath.Ext(fileName) // Извлечение расширения файла

	currentDir, err := os.Getwd() // текущая папка
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	//TODO добавить функцию получения userName из Context
	userName, err := getUserNameFromContext(c)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	var folderToFileStore = "\\storage\\"
	//path := currentDir + folderToFileStore

	fullPathToFile := currentDir + folderToFileStore + userName + "\\" + fileName

	//проверка существования файла в хранилище:
	err = checkFileExist(fullPathToFile)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	userId, err := findUserIdByName(userName)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	//todo определяем размер передаваемого файла / проверка на размер файла доступного пользователю
	fileSizeLimit, err := getFileSizeLimitSql(userName)
	if err != nil {
		//todo log
		return 0, err
	}

	err = checkFileSizeToUpload(fileSizeLimit, c)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	// cохраняем файл в папку
	if err := c.SaveUploadedFile(header, fullPathToFile); err != nil {
		// todo вставить лог
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла", "err: ": err.Error()})
		return 0, err
	}

	//todo функция возвращающая размер создаваемого файла
	fileSize, err := getFileSize(fullPathToFile)
	if err != nil {
		//todo log
		return 0, err
	}

	// to do добавить функционал записи значений в ячейки user_id, file_name, extension, path, в таблицу files
	// TODO AddFileInfoToDB(user_id, file_name, extension, path)
	fileId, err := addFileInfoToDB(userId, fileName, fileExtension, fullPathToFile, fileSize)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	//TODO при сохранении файла создается запись в таблице "accesses", в которой id задачи создается список пользователей,
	//имеющих доступ к данному файлу. При создании файла доступ имеет только создатель.
	err = addAccessInfoToDB(fileId, userId)
	fmt.Println("addAccessInfoToDB: OK")
	if err != nil {
		fmt.Println("addAccessInfoToDB: !OK")
		logger.Error.Println(err.Error())
		return 0, err
	}
	return fileId, nil
}

func (fp *FilePostgres) GetFileByID(fileID int, userName string) (filePath string, err error) {

	userID, err := findUserIdByName(userName)
	if err != nil {
		logger.Error.Println(err.Error())
		return "", err
	}

	err = checkUserToFileAccess(fileID, userID)
	if err != nil {
		logger.Error.Println(err.Error())
		return "", err
	}

	filePath, err = getFilePathByFileID(fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return "", err
	}
	return filePath, err

}

// todo проверить
func (fp *FilePostgres) AllFilesInfo() (files []models.File, err error) {

	rows, err := db.GetDBConn().Query(db.AllFilesInfo)
	if err != nil {
		logger.Error.Println(err.Error())
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
			&file.Added,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		files = append(files, file)
	}

	return files, nil
}

func (fp *FilePostgres) ShowAllUserFilesInfo(c *gin.Context) (files []models.File, err error) {
	userName, err := getUserNameFromContext(c)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}

	userID, err := findUserIdByName(userName)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}

	rows, err := db.GetDBConn().Query(db.GetAllUserFilesSql, userID)
	if err != nil {
		logger.Error.Println(err.Error())
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

// todo доделать функцию (в т.ч. по слоям!)
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

	fmt.Println("Hello from here1")
	row := db.GetDBConn().QueryRow(db.GetFileByFileIDSql, fileID)
	err = row.Scan(
		&file.ID,
		&file.UserId,
		&file.FileName,
		&file.Extension,
		&file.FileSize,
		&file.Added,
	)
	if err != nil {
		logger.Error.Println(err.Error())
		return file, err
	}

	fmt.Println("file: ", file)

	return file, err
}

func (fp *FilePostgres) DeleteFileByID(fileID int) error {

	filePath, err := getFilePathByFileID(fileID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
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

/*
type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres() *TaskPostgres {
	return &TaskPostgres{}
}

func (tp *TaskPostgres) GetAllTasks() (tasks []models.Task, err error) {
	rows, err := db.GetDBConn().Query(db.GetAllTasksSql)
	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := models.Task{}
		err = rows.Scan(
			&t.ID,
			&t.UserId,
			&t.Name,
			&t.Done,
			&t.Description,
			&t.Added,
			&t.DeadLine,
			//&t.DoneAt,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			fmt.Println(err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (tp *TaskPostgres) GetTaskByID(id int) (task models.Task, err error) {
	var t models.Task
	row := db.GetDBConn().QueryRow(db.GetTaskByIDSql, id)

	err = row.Scan(
		&t.ID,
		&t.UserId,
		&t.Name,
		&t.Done,
		&t.Description,
		&t.Added,
		&t.DeadLine,
		//&t.DoneAt,
	)
	if err != nil {
		logger.Error.Println(err.Error())
		return t, err
	}
	return t, nil
}

func (tp *TaskPostgres) GetTaskByUserID(id int) (tasks []models.Task, err error) {
	rows, err := db.GetDBConn().Query(db.GetTaskByUserIDSqlFunc, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := models.Task{}
		err = rows.Scan(
			&t.ID,
			&t.UserId,
			&t.Name,
			&t.Done,
			&t.Description,
			&t.Added,
			&t.DeadLine,
			//&t.DoneAt,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		tasks = append(tasks, t)
	}

	if len(tasks) == 0 {
		return tasks, models.ErrNoRows
	}
	return tasks, nil
}

func (tp *TaskPostgres) GetUndoneTasksByUserID(id int) (tasks []models.Task, err error) {
	rows, err := db.GetDBConn().Query(db.GetUndoneTasksByUserIDSqlFunc, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := models.Task{}
		err = rows.Scan(
			&t.ID,
			&t.UserId,
			&t.Name,
			&t.Done,
			&t.Description,
			&t.Added,
			&t.DeadLine,
			//&t.DoneAt,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		tasks = append(tasks, t)
	}

	if len(tasks) == 0 {
		return tasks, models.ErrNoRows
	}
	return tasks, nil
}

func (tp *TaskPostgres) GetOverdueTasks(id int) (tasks []models.Task, err error) {
	rows, err := db.GetDBConn().Query(db.GetOverdueTasksSqlFunc, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := models.Task{}
		err = rows.Scan(
			&t.ID,
			&t.UserId,
			&t.Name,
			&t.Done,
			&t.Description,
			&t.Added,
			&t.DeadLine,
			//&t.DoneAt,
		)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		tasks = append(tasks, t)
	}
	if len(tasks) == 0 {
		return tasks, models.ErrNoRows
	}
	return tasks, nil
}

func (tp *TaskPostgres) CreateTask(t models.Task) (id int, err error) {
	err = db.GetDBConn().QueryRow(db.CreateTaskSql, t.UserId, t.Name, t.Done, t.Description, t.DeadLine).Scan(&id)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func (tp *TaskPostgres) UpdateTaskByID(id int, t models.Task) (err error) {
	result, err := db.GetDBConn().Exec(db.UpdateTaskByIDSql, t.Name, t.Description, t.Done, t.DeadLine, id)

	updatedRows, _ := result.RowsAffected()
	if updatedRows == 0 {
		return models.ErrNoRows
	}
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}

func (tp *TaskPostgres) ReassignTask(oldUserID, newUserID, id int) error {
	_, err := db.GetDBConn().Exec(db.ReassignTaskSqlProcedure, id, oldUserID, newUserID)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}

func (tp *TaskPostgres) DeleteTaskByID(id int) error {
	result, err := db.GetDBConn().Exec(db.DeleteTaskByIDSql, id)
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
*/
