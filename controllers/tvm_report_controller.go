package controllers

import (
	"backend/models"
	"backend/models/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TVMReportController struct {
	usecase usecase.TVMReportUsecase
}

func NewTVMReportController(usecase usecase.TVMReportUsecase) *TVMReportController {
	return &TVMReportController{usecase: usecase}
}

func (ctrl *TVMReportController) CreateReport(c *gin.Context) {
	var req models.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	report, err := ctrl.usecase.CreateReport(&req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.TvmReportResponse{
		Status:  http.StatusCreated,
		Message: "Report created successfully",
		Data:    report,
	})
}

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

	c.JSON(http.StatusOK, gin.H{"report": report})
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

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"total":   total,
		"page":    filter.Page,
		"limit":   filter.Limit,
	})
}

func (ctrl *TVMReportController) UpdateReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var req models.UpdateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	user, _ := c.Get("user")
	userRole := user.(models.User).Role

	report, err := ctrl.usecase.UpdateReport(uint(id), &req, userID.(uint), userRole)
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
	// user, _ := c.Get("user")

	report, err := ctrl.usecase.UpdateReportStatusByPetugas(
		uint(id),
		req.Status,
		userID.(uint),
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

	ctx.JSON(200, data)
}
