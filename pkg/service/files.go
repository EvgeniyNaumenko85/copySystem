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

//type TaskService struct {
//	repo *repository.Repository
//}
//
//func NewTaskService(repo *repository.Repository) *TaskService {
//	return &TaskService{repo: repo}
//}
//
//func (s *TaskService) GetAllTasks() (tasks []models.Task, err error) {
//	return s.repo.GetAllTasks()
//}
//
//func (s *TaskService) GetTaskByID(id int) (task models.Task, err error) {
//	return s.repo.GetTaskByID(id)
//}
//
//func (s *TaskService) GetOverdueTasks(id int) (tasks []models.Task, err error) {
//	return s.repo.GetOverdueTasks(id)
//}
//
//func (s *TaskService) CreateTask(t models.Task) (id int, err error) {
//	return s.repo.CreateTask(t)
//}
//
//func (s *TaskService) UpdateTaskByID(id int, t models.Task) (err error) {
//	return s.repo.UpdateTaskByID(id, t)
//}
//
//func (s *TaskService) ReassignTask(oldUserID, newUserID, id int) (err error) {
//	return s.repo.ReassignTask(oldUserID, newUserID, id)
//}
//
//func (s *TaskService) DeleteTaskByID(id int) (err error) {
//	return s.repo.DeleteTaskByID(id)
//}
//
//func (s *TaskService) GetTaskByUserID(id int) (tasks []models.Task, err error) {
//	return s.repo.GetTaskByUserID(id)
//}
//
//func (s *TaskService) GetUndoneTasksByUserID(id int) (tasks []models.Task, err error) {
//	return s.repo.GetTaskByUserID(id)
//}
