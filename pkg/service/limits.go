package service

import (
	"copySys/pkg/repository"
)

type LimitsService struct {
	repo *repository.Repository
}

func NewLimitsService(repo *repository.Repository) *LimitsService {
	return &LimitsService{repo: repo}
}

func (ls *LimitsService) SetLimitsToUser(userID, fileSizeLim, storageSizeLim int) (err error) {
	return ls.repo.SetLimitsToUser(userID, fileSizeLim, storageSizeLim)
}
