package repository

import (
	"backend/models"
	"fmt"
	"time"

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
	CountReports() (int64, error)
	CountReportsByStatus() (map[string]int64, error)
	CountReportsPerMonth() ([]map[string]interface{}, error)
	GetLatestReports(limit int) ([]models.TVMReport, error)
	GetOpenReports() ([]models.TVMReport, error)
	TakeReport(reportID uint, teknisiID uint) error
	GetReportsByTechnician(teknisiID uint) ([]models.TVMReport, error)
	ResolveReport(reportID uint, teknisiID uint, note string) error
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
	err := r.db.Preload("Reporter").Preload("Resolver").Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).First(&report, id).Error
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

func (r *tvmReportRepository) CountReports() (int64, error) {
	var total int64
	err := r.db.Model(&models.TVMReport{}).Count(&total).Error
	return total, err
}


func (r *tvmReportRepository) CountReportsByStatus() (map[string]int64, error) {
	type Result struct {
		Status string
		Total  int64
	}

	var results []Result

	err := r.db.
		Model(&models.TVMReport{}).
		Select("status, COUNT(*) as total").
		Group("status").
		Scan(&results).Error

	data := make(map[string]int64)
	for _, r := range results {
		data[r.Status] = r.Total
	}

	return data, err
}

func (r *tvmReportRepository) CountReportsPerMonth() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Raw(`
		SELECT 
			TO_CHAR(created_at, 'Mon') as month,
			COUNT(*) as total
		FROM tvm_reports
		GROUP BY month
		ORDER BY MIN(created_at)
	`).Scan(&results).Error

	return results, err
}

func (r *tvmReportRepository) GetLatestReports(limit int) ([]models.TVMReport, error) {
	var reports []models.TVMReport
	err := r.db.Preload("Reporter").Preload("Resolver").Order("created_at DESC").Limit(limit).Find(&reports).Error

	return reports, err
}

// models/repository/tvm_report_repository.go
func (r *tvmReportRepository) GetOpenReports() ([]models.TVMReport, error) {
	var reports []models.TVMReport
	err := r.db.
		Where("status = ? AND assigned_to IS NULL", models.StatusOpen).
		Order("created_at ASC").
		Find(&reports).Error
	return reports, err
}

func (r *tvmReportRepository) TakeReport(reportID uint, teknisiID uint) error {
	result := r.db.Model(&models.TVMReport{}).
		Where("id = ? AND status = ? AND assigned_to IS NULL",
			reportID, models.StatusOpen).
		Updates(map[string]interface{}{
			"assigned_to": teknisiID,
			"assigned_at": time.Now(),
			"status":      models.StatusInProgress,
		})

	if result.RowsAffected == 0 {
		return fmt.Errorf("laporan sudah diambil atau tidak tersedia")
	}
	return result.Error
}

func (r *tvmReportRepository) GetReportsByTechnician(teknisiID uint) ([]models.TVMReport, error) {
	var reports []models.TVMReport
	query := r.db.Model(&models.TVMReport{}).Preload("Reporter").Preload("Resolver").Preload("Barang").Preload("Assigned")
	err := query.
		Where("assigned_to = ?", teknisiID).
		Order("created_at DESC").
		Find(&reports).Error
	return reports, err
}

// TEKNISI - Resolve Report
func (r *tvmReportRepository) ResolveReport(reportID uint, teknisiID uint, note string) error {
	result := r.db.Model(&models.TVMReport{}).
		Where("id = ? AND status = ? AND assigned_to = ?",
			reportID,
			models.StatusInProgress,
			teknisiID,
		).
		Updates(map[string]interface{}{
			"status":        models.StatusResolved,
			"resolved_by":   teknisiID,
			"resolved_at":   time.Now(),
			"resolved_note": note,
		})

	if result.RowsAffected == 0 {
		return fmt.Errorf("laporan tidak dapat diselesaikan")
	}
	return result.Error
}
