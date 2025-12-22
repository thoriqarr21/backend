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
}

type tvmReportUsecase struct {
	repo repository.TVMReportRepository
}

func NewTVMReportUsecase(repo repository.TVMReportRepository) TVMReportUsecase {
	return &tvmReportUsecase{repo: repo}
}

func (u *tvmReportUsecase) CreateReport(req *models.CreateReportRequest, userID uint) (*models.TVMReport, error) {
	// Check if there's already an active report for this TVM
	existingReport, err := u.repo.FindActiveByTVMCode(req.TVMCode)
	if err == nil && existingReport != nil {
		return nil, errors.New("TVM Code sudah terdaftar. TVM Code: " + req.TVMCode + " dengan status " + string(existingReport.Status))
	}

	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	report := &models.TVMReport{
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

	if req.Priority != "" {
		report.Priority = req.Priority
	}

	if req.Notes != "" {
		report.Notes = req.Notes
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

	return u.repo.Delete(id)
}

func (u *tvmReportUsecase) GetMyReports(userID uint) ([]models.TVMReport, error) {
	return u.repo.GetByUser(userID)
}

func (u *tvmReportUsecase) GetStatistics() (map[string]interface{}, error) {
	return u.repo.GetStatistics()
}