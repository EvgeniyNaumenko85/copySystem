package service

import (
	"copySys/models"
	"copySys/pkg/repository"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password, role string) (string, error)
	ParseToken(token string) (int, string, error)
}

type Task interface {
	GetAllTasks() (tasks []models.Task, err error)
	GetTaskByID(id int) (task models.Task, err error)
	GetOverdueTasks(id int) (tasks []models.Task, err error)
	CreateTask(models.Task) (int, error)
	UpdateTaskByID(id int, t models.Task) (err error)
	ReassignTask(oldUserID, newUserID, id int) (err error)
	DeleteTaskByID(id int) (err error)
	GetTaskByUserID(id int) (tasks []models.Task, err error)
	GetUndoneTasksByUserID(id int) (tasks []models.Task, err error)
}

type File interface {
	//UploadFile(file *multipart.FileHeader, c *gin.Context) (err error)
	UploadFile(file multipart.File, header *multipart.FileHeader, c *gin.Context) (err error)

	//(multipart.File, *multipart.FileHeader, error
	GetFile(id int, c *gin.Context) (err error)
}

type User interface {
	GetAllUsers() (users []models.User, err error)
	GetUserByID(id int) (user models.User, err error)
	//AddUser(u models.User) (id int, err error)
	UpdateUserByID(id int, u models.User) (err error)
	DeleteUserByID(id int) (err error)
}

type Service struct {
	Authorization
	User
	Task
	File
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
		User:          NewUserService(repos),
		Task:          NewTaskService(repos),
		File:          NewFileService(repos),
	}
}
