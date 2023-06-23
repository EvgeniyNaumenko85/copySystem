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
	ParseToken(token string) (int, string, string, error)
}

type Task interface {
	GetAllTasks() (tasks []models.Task, err error)
	GetTaskByID(id int) (task models.Task, err error)
	GetOverdueTasks(id int) (tasks []models.Task, err error)
	CreateTask(models.Task) (int, error)
	UpdateTaskByID(id int, t models.Task) (err error)
	ReassignTask(oldUserID, newUserID, id int) (err error)
	DeleteTaskByID(ID int) (err error)
	GetTaskByUserID(id int) (tasks []models.Task, err error)
	GetUndoneTasksByUserID(id int) (tasks []models.Task, err error)
}

type File interface {
	UploadFile(header *multipart.FileHeader, c *gin.Context) (ID int, err error)
	GetFileByID(fileId int, c *gin.Context) (err error)
	DeleteFileByID(fileID int) (err error)
	ShowAllUserFilesInfo(c *gin.Context) (files []models.File, err error)
	AllFilesInfo() (files []models.File, err error)
}

type User interface {
	GetAllUsers() (users []models.User, err error)
	GetUserByID(id int) (user models.User, err error)
	UpdateUserByID(id int, u models.User) (err error)
	DeleteUserByID(id int) (err error)
}

//(files models.Files )

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
