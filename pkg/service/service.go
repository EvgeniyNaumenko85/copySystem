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

type File interface {
	UploadFile(header *multipart.FileHeader, c *gin.Context) (ID int, err error)
	GetFileByID(fileId int, userName string) (filePath string, err error)
	DeleteFileByID(fileID int) (err error)
	DeleteAllFiles() (err error)
	ShowAllUserFilesInfo(c *gin.Context) (files []models.File, err error)
	AllFilesInfo() (files []models.File, err error)
	FindFileByFileName(fileName, userName string) (file models.File, err error)
}

type User interface {
	GetAllUsers() (users []models.User, err error)
	GetUserByID(id int) (user models.User, err error)
	UpdateUserByID(id int, u models.User) (err error)
	DeleteUserByID(id int) (err error)
}

type Service struct {
	Authorization
	User
	File
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
		User:          NewUserService(repos),
		File:          NewFileService(repos),
	}
}
