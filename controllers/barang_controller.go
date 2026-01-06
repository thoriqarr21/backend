package controllers

import (
	"net/http"
	"strconv"

	"backend/models"
	"backend/models/usecase"

	"github.com/gin-gonic/gin"
)

type BarangController struct {
	usecase usecase.BarangUsecase
}

func NewBarangController(usecase usecase.BarangUsecase) *BarangController {
	return &BarangController{usecase: usecase}
}

func (ctrl *BarangController) GetAllBarang(c *gin.Context) {
	barang, err := ctrl.usecase.GetAllBarang()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"barang": barang})
}

func (ctrl *BarangController) CreateBarang(c *gin.Context) {
	var req models.CreateBarangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		 c.JSON(http.StatusBadRequest, models.ApiResponse{ 
            Status:  http.StatusBadRequest,
            Message: "Invalid request body",
            Data:    err.Error(),
        })
        return
	}

	barang, err := ctrl.usecase.CreateBarang(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
            Status:  http.StatusBadRequest,
            Message: "Failed to create barang",
            Data:    err.Error(),
        })
        return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{
		Status:  http.StatusCreated,
		Message: "Barang created successfully",
		Data:    barang,
	})
}

func (ctrl *BarangController) UpdateBarang(c *gin.Context) {
    // 1️⃣ Validasi ID
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ApiResponse{
            Status:  http.StatusBadRequest,
            Message: "ID harus berupa angka",
        })
        return
    }

    // 2️⃣ Bind & validasi body
    var req models.UpdateBarangRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ApiResponse{ 
            Status:  http.StatusBadRequest,
            Message: "Invalid request body",
            Data:    err.Error(),
        })
        return
    }

    // 3️⃣ Update barang
    barang, err := ctrl.usecase.UpdateBarang(uint(id), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ApiResponse{
            Status:  http.StatusBadRequest,
            Message: "Failed to update barang",
            Data:    err.Error(),
        })
        return
    }

    // 4️⃣ Sukses
    c.JSON(http.StatusOK, models.ApiResponse{
        Status:  http.StatusOK,
        Message: "Barang updated successfully",
        Data:    barang,
    })
}

func (ctrl *BarangController) DeleteBarang(c *gin.Context) {
    // 1️⃣ Validasi ID
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ApiResponse{
            Status:  http.StatusBadRequest,
            Message: "ID harus berupa angka",
        })
        return
    }

    // 2️⃣ Hapus barang
    err = ctrl.usecase.DeleteBarang(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, models.ApiResponse{
            Status:  http.StatusNotFound,
            Message: "Barang not found",
            Data:    err.Error(),
        })
        return
    }

    // 3️⃣ Sukses
    c.JSON(http.StatusOK, models.ApiResponse{
        Status:  http.StatusOK,
        Message: "Barang deleted successfully",
    })
}
