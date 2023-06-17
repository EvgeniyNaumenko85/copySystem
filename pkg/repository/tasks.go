package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"tasks_API/db"
	"tasks_API/models"
	"tasks_API/pkg/logger"
)

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
