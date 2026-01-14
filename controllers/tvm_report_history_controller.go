package controllers

import (
	"backend/models/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TVMReportHistoryController struct {
	usecase usecase.TVMReportHistoryUsecase
}

func NewTVMReportHistoryController(
	usecase usecase.TVMReportHistoryUsecase,
) *TVMReportHistoryController {
	return &TVMReportHistoryController{usecase}
}

func (c *TVMReportHistoryController) GetByReportID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid report id"})
		return
	}

	histories, err := c.usecase.GetByReportID(uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"report_id": id,
		"histories": histories,
	})
}

func (c *TVMReportHistoryController) GetAll(ctx *gin.Context) {
    histories, err := c.usecase.GetAll()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, histories)
}

func (c *TVMReportHistoryController) FindByChangedBy(ctx *gin.Context) {
    userID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    histories, err := c.usecase.GetByChangedBy(uint(userID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "user_id":  userID,
        "histories": histories,
    })
}

func (c *TVMReportHistoryController) GetByUserID(ctx *gin.Context) {
    userID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(400, gin.H{"error": "invalid user id"})
        return
    }

    histories, err := c.usecase.GetByUserID(uint(userID))
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(200, gin.H{
        "user_id":  userID,
        "histories": histories,
    })
}
// func (c *TVMReportHistoryController) GetByID(ctx *gin.Context) {
//     id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
//     if err != nil {
//         ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
//         return
//     }
//     history, err := c.usecase.GetByID(uint(id))
//     if err != nil {
//         ctx.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
//         return
//     }
//     ctx.JSON(http.StatusOK, history)
// }