package usecase

import (
	"backend/models"
	"backend/models/repository"
	"errors"
	"time"
)

type TVMReportUsecase interface {
	CreateReport(req *models.CreateReportRequest, userID uint) (*models.TVMReport, error)
	GetReportByID(id uint) (*models.TVMReport, error)
	GetAllReports(filter *models.ReportFilterRequest) ([]models.TVMReport, int64, error)
	UpdateReport(id uint, req *models.UpdateReportRequest, userID uint, userRole models.Role) (*models.TVMReport, error)
	DeleteReport(id uint, userRole models.Role) error
	GetMyReports(userID uint) ([]models.TVMReport, error)
	GetStatistics() (map[string]interface{}, error)
	UpdateReportStatusByPetugas(
		id uint,
		status models.ReportStatus,
		userID uint,
	) (*models.TVMReport, error)

}

type tvmReportUsecase struct {
	repo    repository.TVMReportRepository
	barangRepo repository.BarangRepository
}

func NewTVMReportUsecase(repo repository.TVMReportRepository, barangRepo repository.BarangRepository) TVMReportUsecase {
	return &tvmReportUsecase{
		repo:      repo,
		barangRepo: barangRepo,
	}
}

func (u *tvmReportUsecase) CreateReport(
	req *models.CreateReportRequest,
	userID uint,
) (*models.TVMReport, error) {

	// 🔍 Cari barang berdasarkan ID
	barang, err := u.barangRepo.FindByID(req.BarangID)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	report := &models.TVMReport{
		BarangID:    barang.ID, // 🔥 DIAMBIL DARI TABEL BARANG
		TVMCode:     req.TVMCode,
		Location:    req.Location,
		IssueType:   req.IssueType,
		Description: req.Description,
		Priority:    priority,
		ImageURL:    req.ImageURL,
		Status:      models.StatusPending,
		ReportedBy:  userID,
	}

	if err := u.repo.Create(report); err != nil {
		return nil, err
	}

	return u.repo.FindByID(report.ID)
}

func (u *tvmReportUsecase) GetReportByID(id uint) (*models.TVMReport, error) {
	return u.repo.FindByID(id)
}

func (u *tvmReportUsecase) GetAllReports(filter *models.ReportFilterRequest) ([]models.TVMReport, int64, error) {
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	return u.repo.GetAll(filter)
}

func (u *tvmReportUsecase) UpdateReport(id uint, req *models.UpdateReportRequest, userID uint, userRole models.Role) (*models.TVMReport, error) {
	report, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("report not found")
	}

	// Only admin or the reporter can update
	if userRole != models.RoleAdmin && report.ReportedBy != userID {
		return nil, errors.New("unauthorized to update this report")
	}

	if req.TVMCode != "" {
		report.TVMCode = req.TVMCode
	}

	// Update fields
	if req.Status != "" {
		report.Status = req.Status
		if req.Status == models.StatusResolved {
			now := time.Now()
			report.ResolvedAt = &now
			if req.ResolvedBy != nil {
				report.ResolvedBy = req.ResolvedBy
			} else {
				report.ResolvedBy = &userID
			}
		}
	}

	if req.IssueType != "" {
		report.IssueType = req.IssueType
	}
	if req.Location != "" {
		report.Location = req.Location
	}
	if req.Description != "" {
		report.Description = req.Description
	}

	if req.Priority != "" {
		report.Priority = req.Priority
	}

	if err := u.repo.Update(report); err != nil {
		return nil, err
	}

	return u.repo.FindByID(report.ID)
}

func (u *tvmReportUsecase) DeleteReport(id uint, userRole models.Role) error {
	// Only admin can delete
	if userRole != models.RoleAdmin {
		return errors.New("only admin can delete reports")
	}

	_, err := u.repo.FindByID(id)
	if err != nil {
		return errors.New("report not found")
	}

	return u.repo.Delete(id)
}

func (u *tvmReportUsecase) GetMyReports(userID uint) ([]models.TVMReport, error) {
	return u.repo.GetByUser(userID)
}

func (u *tvmReportUsecase) GetStatistics() (map[string]interface{}, error) {
	return u.repo.GetStatistics()
}

func (u *tvmReportUsecase) UpdateReportStatusByPetugas(
	id uint,
	status models.ReportStatus,
	userID uint,
) (*models.TVMReport, error) {


	report, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("report not found")
	}

	// 🔁 VALIDASI ALUR STATUS
	switch report.Status {
	case models.StatusPending:
		if status != models.StatusInProgress {
			return nil, errors.New("status tidak valid")
		}

	case models.StatusInProgress:
		if status != models.StatusResolved {
			return nil, errors.New("status tidak valid")
		}

	case models.StatusResolved:
		return nil, errors.New("report sudah selesai")

	default:
		return nil, errors.New("status tidak dikenal")
	}

	// ✅ UPDATE STATUS
	report.Status = status

	// ⏰ JIKA RESOLVED
	if status == models.StatusResolved {
		now := time.Now()
		report.ResolvedAt = &now
		report.ResolvedBy = &userID
	}

	if err := u.repo.Update(report); err != nil {
		return nil, err
	}

	return u.repo.FindByID(report.ID)
}
