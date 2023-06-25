package service

import (
	"copySys/pkg/repository"
)

type AccessService struct {
	repo *repository.Repository
}

func NewAccessService(repo *repository.Repository) *AccessService {
	return &AccessService{repo: repo}
}

func (as *AccessService) ProvidingAccess(fileID, accessToUserID, userID int) (err error) {
	return as.repo.ProvidingAccess(fileID, accessToUserID, userID)
}

func (as *AccessService) ProvidingAccessAll(userID, fileID int) (err error) {
	return as.repo.ProvidingAccessAll(userID, fileID)
}

func (as *AccessService) RemoveAccess(fileID, accessToUserID, userID int) (err error) {
	return as.repo.RemoveAccess(fileID, accessToUserID, userID)
}

func (as *AccessService) RemoveAccessToAll(fileID, userID int) (err error) {
	return as.repo.RemoveAccessToAll(fileID, userID)
}
