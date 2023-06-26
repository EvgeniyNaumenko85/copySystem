package service

import (
	"copySys/models"
	"copySys/pkg/repository"
)

type StatisticsService struct {
	repo *repository.Repository
}

func NewStatisticsService(repo *repository.Repository) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (ss *StatisticsService) GetUserStatistics(userID int) (userStat models.UserStat, err error) {
	return ss.repo.GetUserStatistics(userID)
}
