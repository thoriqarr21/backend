package controllers

import (
	"backend/models"
	"backend/models/usecase"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TVMReportController struct {
	usecase usecase.TVMReportUsecase
}

func NewTVMReportController(usecase usecase.TVMReportUsecase) *TVMReportController {
	return &TVMReportController{usecase: usecase}
}

func (ctrl *TVMReportController) CreateReport(c *gin.Context) {

	var req models.CreateReportRequest

	// 🔴 WAJIB pakai FormMultipart
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 🔍 Ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image wajib diupload",
		})
		return
	}

	// 💾 Simpan gambar
	imagePath, err := usecase.SaveUploadedImage(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 🔐 Ambil user_id dari middleware JWT
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID := userIDAny.(uint)

	// 🚀 Create report
	report, err := ctrl.usecase.CreateReport(
		&req,
		userID,
		imagePath,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  201,
		"message": "Report created successfully",
		"data":    report,
	})
}

func (ctrl *TVMReportController) UpdateReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var req models.UpdateReportRequest
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

		// 🔍 Ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image wajib diupload",
		})
		return
	}

	// 💾 Simpan gambar
	imagePath, err := usecase.SaveUploadedImage(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	user, _ := c.Get("user")
	userRole := user.(models.User).Role

	report, err := ctrl.usecase.UpdateReport(uint(id), &req, models.StatusPending, userID.(uint), userRole, imagePath)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	

	c.JSON(http.StatusOK, models.TvmReportResponse{
		Status:  http.StatusOK,
		Message: "Report updated successfully",
		Data:    report,
	})
}

// func (ctrl *TVMReportController) UpdateReportAdminStatus(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
// 		return
// 	}

// 	var req models.UpdateReportStatusAdminRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	userID, _ := c.Get("user_id")
// 	user, _ := c.Get("user")
// 	userRole := user.(models.User).Role

// 	report, err := ctrl.usecase.UpdateReportAdminStatus(uint(id), &req, userID.(uint), userRole)
// 	if err != nil {
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.TvmReportResponse{
// 		Status:  http.StatusOK,
// 		Message: "Report status updated successfully",
// 		Data:    report,
// 	})
// }

func (ctrl *TVMReportController) GetReportByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	report, err := ctrl.usecase.GetReportByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	baseURL := fmt.Sprintf(
		"%s://%s",
		func() string {
			if c.Request.TLS != nil {
				return "https"
			}
			return "http"
		}(),
		c.Request.Host,
	)

	if report.ImageURL != "" {
		report.ImageURL = baseURL + "/" + report.ImageURL
	}

	c.JSON(http.StatusOK, gin.H{
		"report": report,
	})
}


func (ctrl *TVMReportController) GetAllReports(c *gin.Context) {
	var filter models.ReportFilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reports, total, err := ctrl.usecase.GetAllReports(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	baseURL := fmt.Sprintf(
		"%s://%s",
		func() string {
			if c.Request.TLS != nil {
				return "https"
			}
			return "http"
		}(),
		c.Request.Host,
	)

	for i := range reports {
		if reports[i].ImageURL != "" {
			reports[i].ImageURL = baseURL + "/" + reports[i].ImageURL
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"total":   total,
		"page":    filter.Page,
		"limit":   filter.Limit,
	})
}

func (ctrl *TVMReportController) DeleteReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	user, _ := c.Get("user")
	userRole := user.(models.User).Role

	if err := ctrl.usecase.DeleteReport(uint(id), userRole); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.TvmReportResponse{
		Status:  http.StatusOK,
		Message: "Report deleted successfully",
	})
}

func (ctrl *TVMReportController) GetMyReports(c *gin.Context) {
	userID, _ := c.Get("user_id")

	reports, err := ctrl.usecase.GetMyReports(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
		baseURL := fmt.Sprintf(
		"%s://%s",
		func() string {
			if c.Request.TLS != nil {
				return "https"
			}
			return "http"
		}(),
		c.Request.Host,
	)

	for i := range reports {
		if reports[i].ImageURL != "" {
			reports[i].ImageURL = baseURL + "/" + reports[i].ImageURL
		}
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

func (ctrl *TVMReportController) GetStatistics(c *gin.Context) {
	stats, err := ctrl.usecase.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statistics": stats})
}

func (ctrl *TVMReportController) UpdateStatusByPetugas(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var req models.UpdateReportStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	roleVal, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
		return
	}

	userRole := models.Role(roleVal.(string))

	report, err := ctrl.usecase.UpdateReportStatusByPetugas(
		uint(id),
		req.Status,
		userID.(uint),
		userRole,
	)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.TvmReportResponse{
		Status:  http.StatusOK,
		Message: "Status berhasil diperbarui",
		Data:    report,
	})
}

func (c *TVMReportController) GetDashboard(ctx *gin.Context) {
    userID, _ := ctx.Get("user_id")
    data, err := c.usecase.GetDashboard(userID.(uint))
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }

    baseURL := fmt.Sprintf(
        "%s://%s",
        func() string {
            if ctx.Request.TLS != nil {
                return "https"
            }
            return "http"
        }(),
        ctx.Request.Host,
    )

    // Process latest_reports to update ImageURL if it exists
    if latestReports, ok := data["latest_reports"].([]models.TVMReport); ok {
        for i := range latestReports {
            if latestReports[i].ImageURL != "" {
                latestReports[i].ImageURL = baseURL + "/" + latestReports[i].ImageURL
            }
        }
        data["latest_reports"] = latestReports
    }

    ctx.JSON(200, data)
}

// controllers/tvm_report_controller.go

// 🔓 List laporan PENDING
func (c *TVMReportController) GetPendingReports(ctx *gin.Context) {
	reports, err := c.usecase.GetPendingReports()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reports)
}

// 🔒 Ambil laporan
func (c *TVMReportController) TakeReport(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status := models.StatusInProgress 
	err := c.usecase.TakeReport(
		uint(id),
		userID.(uint),
		status,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "laporan berhasil diambil"})
}

func (c *TVMReportController) ResolveReport(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status := models.StatusResolved
	err := c.usecase.ResolveReport(
		uint(id),
		userID.(uint),
		status,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "laporan berhasil diambil"})
}

func (c *TVMReportController) OpenReport(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.usecase.OpenReport(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "laporan berhasil dibuka",
	})
}

func (ctrl *TVMReportController) GetMyReportsAsTechnician(c *gin.Context) {
	userID, _ := c.Get("user_id")

	reports, err := ctrl.usecase.GetMyReportsAsTechnician(userID.(uint))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"reports": reports})
}

// func (c *TVMReportController) ResolveReport(ctx *gin.Context) {
// 	id, _ := strconv.Atoi(ctx.Param("id"))
// 	teknisiID, _ := ctx.Get("user_id")

// 	var req struct {
// 		Note string `json:"note"`
// 	}

// 	ctx.ShouldBindJSON(&req)

// 	err := c.usecase.ResolveReport(uint(id), teknisiID.(uint), req.Note)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "laporan berhasil diselesaikan"})
// }
