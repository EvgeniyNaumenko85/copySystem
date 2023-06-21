package repository

import (
	"copySys/db"
	"copySys/models"
	"copySys/pkg/logger"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FilePostgres struct {
	db *sqlx.DB
}

func NewFilePostgres() *FilePostgres {
	return &FilePostgres{}
}

func fileExists(filename string) error {
	_, err := os.Stat(filename)
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
		return "", errors.New("can't  get userName")
	} else {
		userName := fmt.Sprintf("%v", userNameTypeAny)
		return userName, nil
	}
}

func findUserIdByName(userName string) (int, error) {
	//fmt.Println("userName: ", userName)
	if userName == "" {
		return -1, models.ErrUserNotExists
	}

	var id int
	err := db.GetDBConn().QueryRow(db.GetIdUserByNameSql, userName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, models.ErrUserNotExists
		} else {
			logger.Error.Println(err.Error())
			fmt.Println(err.Error())
			return -1, err
		}
	}

	return id, nil
}

func fileSizeToUpload(fileName string, c *gin.Context) (fileSize int, err error) {
	file, err := c.FormFile(fileName)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}

	fileSizeInBytes := file.Size
	fileSize = int(math.Round(float64(fileSizeInBytes) / (1024 * 1024)))

	return fileSize, nil
}

func checkFileSizeLimitSql(userName string, fileSize int) error {
	var fileSizeLimit int
	err := db.GetDBConn().QueryRow(db.CheckFileSizeLimitSql, userName).Scan(&fileSizeLimit)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	if fileSize > fileSizeLimit {
		return models.FileToBig
	}

	return nil
}

func addFileInfoToDB(userId int, fileName, extension, path string) (id int, err error) {
	err = db.GetDBConn().QueryRow(db.CreateFileSql, userId, fileName, extension, path).Scan(&id)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func addAccessInfoToDB(fileId, userId int) error {
	_, err := db.GetDBConn().Exec(db.CreateAccessSql, fileId, userId)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}

func userToFileAccess(fileId, userId int) error {
	result, err := db.GetDBConn().Exec(db.CheckAccessInTableSql, fileId, userId)
	foundRows, _ := result.RowsAffected()
	if foundRows == 0 {
		return models.FileAccessDenied
	}
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}

	return nil
}

func (fp *FilePostgres) UploadFile(file multipart.File, header *multipart.FileHeader, c *gin.Context) (int, error) {
	fileName := filepath.Base(header.Filename)
	extension := filepath.Ext(fileName) // Извлечение расширения файла
	//onlyName := filename[:len(filename)-len(extension)] // Извлечение имени файла без расширения

	//destinationPath := `C:\Users\Евгений Науменко\Desktop\copySys\storage\` + filename //  временная затычка.
	//destinationPath := `.\copySys\storage\` + filename //  временная затычка
	currentDir, err := os.Getwd() // текущая папка
	if err != nil {
		//todo log.err
		return 0, err
	}

	path := currentDir + "\\storage\\"

	//TODO добавить функцию получения userName из Context
	userName, err := getUserNameFromContext(c)
	//fmt.Println("userName, err :", userName, err)
	if err != nil {
		//todo log
		return 0, err
	}

	toFileExists := currentDir + "\\storage\\" + userName + "\\" + "\\" + fileName

	//проверка существования файла в хранилище:
	//todo не забыть включить эту проверку
	err = fileExists(toFileExists)
	if err != nil {
		//todo log
		return 0, err
	}

	userId, err := findUserIdByName(userName)
	if err != nil {
		//todo log
		return 0, err
	}

	//todo определяем размер передаваемого файла
	fileSize, err := fileSizeToUpload(fileName, c)
	if err != nil {
		//todo log
		fmt.Println("HIIII")
		fmt.Println("fileName:", fileName)
		fmt.Println("fileSize", fileSize)
		return 0, err
	}

	//todo проверяем размер загружаемого файла, по user_id смотрим чтоб fileSize <= size_limit в таблице users
	err = checkFileSizeLimitSql(userName, fileSize)
	if err != nil {
		//todo log
		return 0, err
	}

	// cохраняем файл в папку
	if err := c.SaveUploadedFile(header, toFileExists); err != nil {
		// todo вставить лог
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла", "err: ": err.Error()})
		return 0, err
	}

	// to do добавить функционал записи значений в ячейки user_id, file_name, extension, path, в таблицу files
	// TODO AddFileInfoToDB(user_id, file_name, extension, path)

	fileId, err := addFileInfoToDB(userId, fileName, extension, path)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	//TODO при сохранении файла создается запись в таблице "accesses",в которой id задачи создается список пользователей,
	//имеющих доступ к данному файлу. При создании файла доступ имеет только создатель.

	err = addAccessInfoToDB(fileId, userId)

	return fileId, nil
}

func (fp *FilePostgres) GetFile(fileId int, c *gin.Context) (err error) {

	//TODO добавить функцию получения userName из Context
	userName, err := getUserNameFromContext(c)
	//fmt.Println("userName, err :", userName, err)
	if err != nil {
		//todo log
		return err
	}

	userId, err := findUserIdByName(userName)
	if err != nil {
		//todo log
		return err
	}

	//todo функция прав доступа к файлу (ищет пару user_id/file_id в таблице access)
	err = userToFileAccess(fileId, userId)
	if err != nil {
		//todo log
		return err
	}

	//
	fmt.Println("Hello from postgres.loadFile")
	//получить ссылку на файл path из таблицы files  по переданному id
	//path := `C:\Users\Евгений Науменко\Desktop\copySys\storage\testFile.txt` //искусственная "заглушка"
	path := `C:\Users\Евгений Науменко\Desktop\copySys\storage\test.exe`
	//path := `C:\Users\Евгений Науменко\Desktop\copySys\storage\test.txt`

	//filename := "testFile.txt" // берем из таблицы files по id

	fmt.Println(filepath.Base(path))

	// Устанавливаем заголовки для скачивания файла
	c.Header("Content-Description", "File Transfer")
	//c.Header("Content-Disposition", "attachment; filename="+filename)
	//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	//c.Header("Content-Disposition", "attachment; filename="+filename)
	//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
	//c.File(path) //ругается на  http: wrote more than the declared Content-Length

	// Устанавливаем заголовок Content-Disposition для передачи имени файла с расширением
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))

	// Открываем файл для чтения
	file, err := os.Open(path)
	if err != nil {

		return err
	}
	defer file.Close()

	// Отправляем файл в ответ
	_, err = io.Copy(c.Writer, file)
	if err != nil {

		return err
	}

	return
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
