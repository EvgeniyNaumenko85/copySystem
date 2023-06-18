package service

import (
	"copySys/models"
	"copySys/pkg/repository"
)

type TaskService struct {
	repo *repository.Repository
}

func NewTaskService(repo *repository.Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAllTasks() (tasks []models.Task, err error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) GetTaskByID(id int) (task models.Task, err error) {
	return s.repo.GetTaskByID(id)
}

func (s *TaskService) GetOverdueTasks(id int) (tasks []models.Task, err error) {
	return s.repo.GetOverdueTasks(id)
}

func (s *TaskService) CreateTask(t models.Task) (id int, err error) {
	return s.repo.CreateTask(t)
}

func (s *TaskService) UpdateTaskByID(id int, t models.Task) (err error) {
	return s.repo.UpdateTaskByID(id, t)
}

func (s *TaskService) ReassignTask(oldUserID, newUserID, id int) (err error) {
	return s.repo.ReassignTask(oldUserID, newUserID, id)
}

func (s *TaskService) DeleteTaskByID(id int) (err error) {
	return s.repo.DeleteTaskByID(id)
}

func (s *TaskService) GetTaskByUserID(id int) (tasks []models.Task, err error) {
	return s.repo.GetTaskByUserID(id)
}

func (s *TaskService) GetUndoneTasksByUserID(id int) (tasks []models.Task, err error) {
	return s.repo.GetTaskByUserID(id)
}
