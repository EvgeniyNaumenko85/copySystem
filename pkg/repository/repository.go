package repository

import (
	"copySys/models"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password, role string) (models.User, error)
}

type User interface {
	GetAllUsers() (users []models.User, err error)
	GetUserByID(id int) (user models.User, err error)
	UpdateUserByID(id int, u models.User) (err error)
	DeleteUserByID(id int) (err error)
}

type File interface {
	UploadFile(header *multipart.FileHeader, c *gin.Context) (id int, err error)
	GetFileByID(fileID int, userName string) (filePath string, err error)
	DeleteFileByID(fileID int) (err error)
	DeleteAllFiles() (err error)
	ShowAllUserFilesInfo(c *gin.Context) (files []models.File, err error)
	AllFilesInfo() (files []models.File, err error)
	FindFileByFileName(fileName, userName string) (file models.File, err error)
}

type Repository struct {
	Authorization
	User
	File
}

func NewRepository() *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(),
		User:          NewUserPostgres(),
		File:          NewFilePostgres(),
	}
}
