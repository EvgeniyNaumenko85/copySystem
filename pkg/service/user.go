package service

import (
	"copySys/models"
	"copySys/pkg/repository"
	"fmt"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() (users []models.User, err error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) GetUserByID(id int) (user models.User, err error) {
	return s.repo.GetUserByID(id)
}

//func (s *UserService) CreateUser(u models.User) (id int, err error) {
//	return s.repo.CreateUser(u)
//}

func (s *UserService) UpdateUserByID(id int, u models.User) (err error) {
	fmt.Println("UpdateUserByID:", id, u)
	return s.repo.UpdateUserByID(id, u)
}

func (s *UserService) DeleteUserByID(id int) (err error) {
	return s.repo.DeleteUserByID(id)
}

//func (s *UserService) CreateUser(u models.User) (id int, err error) {
//	u.Password = generatePasswordHash(u.Password)
//	return s.repo.AddUser(u)
//}
