package usecase

import (
	"backend/models"
	"backend/models/repository"
)

type TVMReportHistoryUsecase interface {
	GetByReportID(reportID uint) ([]models.TVMReportHistory, error)
	GetByUserID(userID uint) ([]models.TVMReportHistory, error)
	GetAll() ([]models.TVMReportHistory, error)
	GetByChangedBy(userID uint) ([]models.TVMReportHistory, error)
}

type tvmReportHistoryUsecase struct {
	repo repository.TVMReportHistoryRepository
}

func NewTVMReportHistoryUsecase(
	repo repository.TVMReportHistoryRepository,
) TVMReportHistoryUsecase {
	return &tvmReportHistoryUsecase{repo}
}

func (uc *tvmReportHistoryUsecase) GetByReportID(
	reportID uint,
) ([]models.TVMReportHistory, error) {
	return uc.repo.FindByReportID(reportID)
}

func (uc *tvmReportHistoryUsecase) GetByUserID(userID uint) ([]models.TVMReportHistory, error) {
	return uc.repo.FindByUserID(userID)
}

func (uc *tvmReportHistoryUsecase) GetAll() ([]models.TVMReportHistory, error) {
	return uc.repo.GetAll()
}

func (uc *tvmReportHistoryUsecase) GetByChangedBy(userID uint) ([]models.TVMReportHistory, error) {
	return uc.repo.FindByChangedBy(userID)
}
