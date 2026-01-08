package repository

import (
	"backend/models"

	"gorm.io/gorm"
)

type TVMReportRepository interface {
	Create(report *models.TVMReport) error
	FindByID(id uint) (*models.TVMReport, error)
	GetAll(filter *models.ReportFilterRequest) ([]models.TVMReport, int64, error)
	Update(report *models.TVMReport) error
	Delete(id uint) error
	GetByUser(userID uint) ([]models.TVMReport, error)
	GetStatistics() (map[string]interface{}, error)
	FindActiveByTVMCode(tvmCode string) (*models.TVMReport, error)
	UpdateStatus(id uint, status models.ReportStatus, resolverID *uint) error

}

type tvmReportRepository struct {
	db *gorm.DB
}

// FindActiveByTVMCode implements [TVMReportRepository].
func (r *tvmReportRepository) FindActiveByTVMCode(tvmCode string) (*models.TVMReport, error) {
	var report models.TVMReport
	err := r.db.Where("tvm_code = ? AND status IN (?)", 
		tvmCode, 
		[]string{string(models.StatusPending), string(models.StatusInProgress)},
	).First(&report).Error
	
	if err != nil {
		return nil, err
	}
	return &report, nil
}
func NewTVMReportRepository(db *gorm.DB) TVMReportRepository {
	return &tvmReportRepository{db: db}
}

func (r *tvmReportRepository) Create(report *models.TVMReport) error {
	return r.db.Create(report).Error
}

func (r *tvmReportRepository) FindByID(id uint) (*models.TVMReport, error) {
	var report models.TVMReport
	err := r.db.Preload("Reporter").Preload("Resolver").First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *tvmReportRepository) GetAll(filter *models.ReportFilterRequest) ([]models.TVMReport, int64, error) {
	var reports []models.TVMReport
	var total int64

	query := r.db.Model(&models.TVMReport{}).Preload("Reporter").Preload("Resolver").Preload("Barang")

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.TVMCode != "" {
		query = query.Where("tvm_code LIKE ?", "%"+filter.TVMCode+"%")
	}
	if filter.ReportedBy != 0 {
		query = query.Where("reported_by = ?", filter.ReportedBy)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	err := query.Order("created_at DESC").Find(&reports).Error
	return reports, total, err
}

func (r *tvmReportRepository) Update(report *models.TVMReport) error {
	return r.db.Save(report).Error
}

func (r *tvmReportRepository) Delete(id uint) error {
	return r.db.Delete(&models.TVMReport{}, id).Error
}

func (r *tvmReportRepository) GetByUser(userID uint) ([]models.TVMReport, error) {
	var reports []models.TVMReport
	err := r.db.Where("reported_by = ?", userID).
		Preload("Reporter").
		Preload("Resolver").
		Order("created_at DESC").
		Find(&reports).Error
	return reports, err
}

func (r *tvmReportRepository) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	r.db.Model(&models.TVMReport{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts)

	stats["by_status"] = statusCounts

	// Count by priority
	var priorityCounts []struct {
		Priority string
		Count    int64
	}
	r.db.Model(&models.TVMReport{}).
		Select("priority, count(*) as count").
		Group("priority").
		Scan(&priorityCounts)

	stats["by_priority"] = priorityCounts

	// Total reports
	var total int64
	r.db.Model(&models.TVMReport{}).Count(&total)
	stats["total"] = total

	return stats, nil
}

func (r *tvmReportRepository) UpdateStatus(
	id uint,
	status models.ReportStatus,
	resolverID *uint,
) error {
	updateData := map[string]interface{}{
		"status": status,
	}

	// Jika ada resolver (petugas)
	if resolverID != nil {
		updateData["resolved_by"] = *resolverID
	}

	return r.db.Model(&models.TVMReport{}).
		Where("id = ?", id).
		Updates(updateData).Error
}
