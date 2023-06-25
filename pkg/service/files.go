package service

import (
	"copySys/models"
	"copySys/pkg/repository"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type FileService struct {
	repo *repository.Repository
}

func NewFileService(repo *repository.Repository) *FileService {
	return &FileService{repo: repo}
}

func (fs *FileService) UploadFile(header *multipart.FileHeader, c *gin.Context) (fileId int, err error) {
	return fs.repo.UploadFile(header, c)
}

func (fs *FileService) GetFileByID(fileID int, userName string) (filePath string, err error) {
	return fs.repo.GetFileByID(fileID, userName)
}

func (fs *FileService) AllFilesInfo() (files []models.File, err error) {
	return fs.repo.AllFilesInfo()
}

func (fs *FileService) ShowAllUserFilesInfo(c *gin.Context) (files []models.File, err error) {
	return fs.repo.ShowAllUserFilesInfo(c)
}

func (fs *FileService) FindFileByFileName(fileName, userName string) (file models.File, err error) {
	return fs.repo.FindFileByFileName(fileName, userName)
}

func (fs *FileService) DeleteFileByID(fileID int) (err error) {
	return fs.repo.DeleteFileByID(fileID)
}
