package repository

import (
	"backend/models"

	"gorm.io/gorm"
)

type TVMReportHistoryRepository interface {
	Create(history *models.TVMReportHistory) error
	FindByReportID(reportID uint) ([]models.TVMReportHistory, error)
	FindByUserID(userID uint) ([]models.TVMReportHistory, error)
	GetAll() ([]models.TVMReportHistory, error)
	FindByChangedBy(userID uint) ([]models.TVMReportHistory, error)
}

type tvmReportHistoryRepository struct {
	db *gorm.DB
}

func NewTVMReportHistoryRepository(db *gorm.DB) TVMReportHistoryRepository {
	return &tvmReportHistoryRepository{db}
}

func (r *tvmReportHistoryRepository) Create(h *models.TVMReportHistory) error {
	return r.db.Create(h).Error
}

func (r *tvmReportHistoryRepository) FindByReportID(id uint) ([]models.TVMReportHistory, error) {
	var histories []models.TVMReportHistory
	query := r.db.Model(&models.TVMReportHistory{}).Preload("Report").Preload("ChangedByUser")
	err := query.
		Where("tvm_report_id = ?", id).
		Order("created_at ASC").
		Find(&histories).Error
	return histories, err
}

func (r *tvmReportHistoryRepository) FindByUserID(userID uint) ([]models.TVMReportHistory, error) {
    var histories []models.TVMReportHistory
    query := r.db.Model(&models.TVMReportHistory{}).Preload("Report").Preload("ChangedByUser")
    err := query.Where("user_id = ?", userID).Find(&histories).Error
    if err != nil {
        return nil, err
    }
    return histories, nil
}

func (r *tvmReportHistoryRepository) FindByChangedBy(userID uint) ([]models.TVMReportHistory, error) {
    var histories []models.TVMReportHistory
	query := r.db.Model(&models.TVMReportHistory{}).Preload("ChangedByUser")
    err := query.
        Where("changed_by = ?", userID).
        Order("created_at DESC").
        Find(&histories).Error
    return histories, err
}

func (r *tvmReportHistoryRepository) GetAll() ([]models.TVMReportHistory, error) {
    var histories []models.TVMReportHistory
	query := r.db.Model(&models.TVMReportHistory{}).Preload("ChangedByUser")
    err := query.Find(&histories).Error
    if err != nil {
        return nil, err
    }
    return histories, nil
}