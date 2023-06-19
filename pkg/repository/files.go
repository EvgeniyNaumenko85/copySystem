package repository

import (
	"copySys/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"io"
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
		fmt.Println("err: ", err)
		return models.ErrFileAlredyExists
	}
}

// func (fp *FilePostgres) UploadFile(file *multipart.FileHeader, c *gin.Context) error {
func (fp *FilePostgres) UploadFile(file multipart.File, header *multipart.FileHeader, c *gin.Context) error {
	filename := filepath.Base(header.Filename)
	//destinationPath := `C:\Users\Евгений Науменко\Desktop\copySys\storage\` + filename //  временная затычка.
	destinationPath := `.\copySys\storage\` + filename //  временная затычка
	currentDir, err := os.Getwd()                      // текущая папка
	if err != nil {
		// Обработка ошибки
	}
	fmt.Println("currentDir: ", currentDir)
	fmt.Println(destinationPath)
	//err = fileExists(destinationPath)

	tofileExists := currentDir + "\\storage\\" + filename
	fmt.Println("tofileExists: ", tofileExists)
	err = fileExists(tofileExists)
	if err != nil {
		return err
	} else {
		if err != nil {
			return err
		}
	}

	if err := c.SaveUploadedFile(header, tofileExists); err != nil {
		//вставить лог
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла", "err: ": err.Error()})
		return err
	}

	// to do добавить функционал записи значений в ячейки user_id, file_name, extension, path,  в таблицу files

	return nil
}

func (fp *FilePostgres) GetFile(id int, c *gin.Context) (err error) {
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
