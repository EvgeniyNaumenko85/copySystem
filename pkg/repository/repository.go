package repository

import (
	"tasks_API/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password, role string) (models.User, error)
}

type Task interface {
	GetAllTasks() (tasks []models.Task, err error)
	GetTaskByID(id int) (task models.Task, err error)
	GetOverdueTasks(id int) (tasks []models.Task, err error)
	CreateTask(models.Task) (id int, err error)
	UpdateTaskByID(id int, t models.Task) (err error)
	ReassignTask(oldUserID, newUserID, id int) (err error)
	DeleteTaskByID(id int) (err error)
	GetTaskByUserID(id int) (tasks []models.Task, err error)
	GetUndoneTasksByUserID(id int) (tasks []models.Task, err error)
}

type User interface {
	GetAllUsers() (users []models.User, err error)
	GetUserByID(id int) (user models.User, err error)
	UpdateUserByID(id int, u models.User) (err error)
	DeleteUserByID(id int) (err error)
}

type Repository struct {
	Authorization
	User
	Task
}

func NewRepository() *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(),
		User:          NewUserPostgres(),
		Task:          NewTaskPostgres(),
	}
}

type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
	SSLMode  string
}
