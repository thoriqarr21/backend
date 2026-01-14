package usecase

import (
	"backend/models"
	"backend/models/repository"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type TVMReportUsecase interface {
	CreateReport(req *models.CreateReportRequest, userID uint, imagePath string) (*models.TVMReport, error)
	GetReportByID(id uint) (*models.TVMReport, error)
	GetAllReports(filter *models.ReportFilterRequest) ([]models.TVMReport, int64, error)
	UpdateReport(id uint, req *models.UpdateReportRequest, userID uint, userRole models.Role, imagePath string) (*models.TVMReport, error)
	DeleteReport(id uint, userRole models.Role) error
	GetMyReports(userID uint) ([]models.TVMReport, error)
	GetDashboard(userID uint) (map[string]interface{}, error)
	GetStatistics() (map[string]interface{}, error)
	UpdateReportStatusByPetugas(
		id uint,
		status models.ReportStatus,
		userID uint,
		userRole models.Role,
	) (*models.TVMReport, error)
	TakeReport(reportID uint, teknisiID uint) error
	GetOpenReports() ([]models.TVMReport, error)
	OpenReport(reportID uint) error
	GetMyReportsAsTechnician(userID uint) ([]models.TVMReport, error)
	ResolveReport(reportID uint, teknisiID uint, note string) error
}

type tvmReportUsecase struct {
	repo    repository.TVMReportRepository
	barangRepo repository.BarangRepository
	historyRepo repository.TVMReportHistoryRepository
}

func NewTVMReportUsecase(repo repository.TVMReportRepository, barangRepo repository.BarangRepository, historyRepo repository.TVMReportHistoryRepository) TVMReportUsecase {
	return &tvmReportUsecase{
		repo:      repo,
		barangRepo: barangRepo,
		historyRepo: historyRepo,
	}
}

func (u *tvmReportUsecase) CreateReport(
    req *models.CreateReportRequest,
    userID uint,
    imagePath string, // Ini akan berisi "uploads/nama-file.jpg"
) (*models.TVMReport, error) {

    barang, err := u.barangRepo.FindByID(req.BarangID)
    if err != nil {
        return nil, errors.New("barang tidak ditemukan")
    }

    stok, _ := strconv.Atoi(barang.Stok)
    if stok <= 0 {
        return nil, errors.New("stok barang habis")
    }

    barang.Stok = strconv.Itoa(stok - 1)
    _ = u.barangRepo.Update(barang)

    report := &models.TVMReport{
        BarangID:    barang.ID,
        TVMCode:     req.TVMCode,
        Location:    req.Location,
        IssueType:   req.IssueType,
        Description: req.Description,
        Priority:    req.Priority,
        ImageURL:    imagePath, // Menyimpan path lokal ke database
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

func (u *tvmReportUsecase) UpdateReport(
	id uint,
	req *models.UpdateReportRequest,
	userID uint,
	userRole models.Role,
	imagePath string, // hasil upload baru (opsional)
) (*models.TVMReport, error) {

	report, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("report not found")
	}

	// 🔐 Authorization
	if userRole != models.RoleAdmin && report.ReportedBy != userID {
		return nil, errors.New("unauthorized to update this report")
	}

	// 📝 Update fields (partial)
	if req.TVMCode != "" {
		report.TVMCode = req.TVMCode
	}

	if req.Status != "" {
		report.Status = req.Status

		if req.Status == models.StatusResolved {
			now := time.Now()
			report.ResolvedAt = &now

			if userRole == models.RoleAdmin {
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

	// 🖼️ Update image (HANYA dari upload)
	if imagePath != "" {
		report.ImageURL = imagePath
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
	userRole models.Role,
) (*models.TVMReport, error) {

	report, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("report not found")
	}

	// 🔴 SIMPAN STATUS LAMA UNTUK HISTORY
	oldStatus := report.Status

	switch userRole {

	// 🧑‍💼 PETUGAS
	case models.RoleUser:
		if !(report.Status == models.StatusPending && status == models.StatusOpen) {
			return nil, fmt.Errorf("invalid status transition for petugas")
		}
		report.Status = status

	// 🧑‍🔧 TEKNISI
	case models.RoleTeknisi:

		// 📌 Ambil laporan
		if report.Status == models.StatusOpen && status == models.StatusInProgress {

			if report.AssignedTo != nil {
				return nil, fmt.Errorf("laporan sudah diambil teknisi lain")
			}

			now := time.Now()
			report.Status = models.StatusInProgress
			report.AssignedTo = &userID
			report.AssignedAt = &now

		// 📌 Selesaikan laporan
		} else if report.Status == models.StatusInProgress && status == models.StatusResolved {

			if report.AssignedTo == nil || *report.AssignedTo != userID {
				return nil, fmt.Errorf("anda bukan teknisi yang menangani laporan ini")
			}

			now := time.Now()
			report.Status = models.StatusResolved
			report.ResolvedAt = &now
			report.ResolvedBy = &userID

		} else {
			return nil, fmt.Errorf("invalid status transition for teknisi")
		}

	// 🧑‍💼 ADMIN
	case models.RoleAdmin:
		if status != models.StatusRejected && status != models.StatusOpen {
			return nil, fmt.Errorf("invalid status transition for admin")
		}
		report.Status = status

	default:
		return nil, fmt.Errorf("invalid role")
	}

	// ✅ UPDATE REPORT
	if err := u.repo.Update(report); err != nil {
		return nil, err
	}

	// 🕒 SIMPAN HISTORY
	history := &models.TVMReportHistory{
		TVMReportID: report.ID,
		FromStatus:  oldStatus,
		ToStatus:    report.Status,
		ChangedBy:   userID,
		Role:        userRole,
	}

	_ = u.historyRepo.Create(history)

	return u.repo.FindByID(report.ID)
}



func (u *tvmReportUsecase) GetDashboard(userID uint) (map[string]interface{}, error) {

	totalReports, err := u.repo.CountReports()
	if err != nil {
		return nil, err
	}

	totalBarang, err := u.barangRepo.CountBarang()
	if err != nil {
		return nil, err
	}

	statusStats, err := u.repo.CountReportsByStatus()
	if err != nil {
		return nil, err
	}

	monthlyStats, err := u.repo.CountReportsPerMonth()
	if err != nil {
		return nil, err
	}

	latestReports, err := u.repo.GetLatestReports(5)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"total_reports":     totalReports,
		"total_barang":      totalBarang,
		"reports_by_status": statusStats,
		"reports_per_month": monthlyStats,
		"latest_reports":    latestReports,
	}, nil
}

func SaveUploadedImage(ctx *gin.Context, file *multipart.FileHeader) (string, error) {

	// ✅ PASTIKAN FOLDER ADA
	if err := os.MkdirAll("uploads/reports", 0755); err != nil {
		return "", err
	}

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", errors.New("format gambar harus jpg, jpeg, atau png")
	}

	// Nama file unik
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	path := "uploads/reports/" + filename

	// ✅ SIMPAN FILE
	if err := ctx.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	return path, nil
}

func (u *tvmReportUsecase) GetOpenReports() ([]models.TVMReport, error) {
	return u.repo.GetOpenReports()
}

func (u *tvmReportUsecase) TakeReport(reportID uint, teknisiID uint) error {
	return u.repo.TakeReport(reportID, teknisiID)
}

func (u *tvmReportUsecase) OpenReport(reportID uint) error {
	report, err := u.repo.FindByID(reportID)
	if err != nil {
		return errors.New("laporan tidak ditemukan")
	}

	if report.Status != models.StatusPending {
		return errors.New("laporan tidak bisa dibuka")
	}

	report.Status = models.StatusOpen
	return u.repo.Update(report)
}

func (u *tvmReportUsecase) GetMyReportsAsTechnician(teknisiID uint) ([]models.TVMReport, error) {
	return u.repo.GetReportsByTechnician(teknisiID)
}

func (u *tvmReportUsecase) ResolveReport(
	reportID uint,
	teknisiID uint,
	note string,
) error {
	return u.repo.ResolveReport(reportID, teknisiID, note)
}



