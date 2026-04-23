package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mobile-api/models"
	"mobile-api/utils"
)

type BookingController struct {
	db *gorm.DB
}

func NewBookingController(db *gorm.DB) *BookingController {
	return &BookingController{db: db}
}

// GET /api/v1/bookings
func (bc *BookingController) GetMyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status   := c.Query("status")

	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }

	// Join t_flight_transactions → t_contacts via user phone/email
	query := bc.db.Model(&models.FlightBooking{}).
		Joins("JOIN t_flight_transactions ft ON ft.id = t_flight_bookings.transaction_id").
		Joins("JOIN t_contacts tc ON tc.id = ft.contact_id").
		Where("ft.id IN (SELECT id FROM t_flight_transactions WHERE contact_id IN (SELECT id FROM t_contacts WHERE email = (SELECT email FROM users WHERE id = ?)))", userID)

	if status != "" {
		query = query.Where("t_flight_bookings.status = ?", status)
	}

	var total int64
	query.Count(&total)

	var bookings []models.FlightBooking
	offset := (page - 1) * limit
	if err := query.
		Preload("Passengers").
		Preload("Segments").
		Preload("Remarks").
		Order("t_flight_bookings.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch bookings", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Bookings retrieved", gin.H{
		"bookings": bookings,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_page": (total + int64(limit) - 1) / int64(limit),
		},
	}))
}

// GET /api/v1/bookings/:id
func (bc *BookingController) GetBookingDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID", ""))
		return
	}

	var booking models.FlightBooking
	if err := bc.db.
		Preload("Passengers").
		Preload("Segments").
		Preload("Remarks").
		First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
		return
	}

	// Load transaction & contact
	var tx models.FlightTransaction
	bc.db.Preload("Contact").Preload("PaymentMethod").First(&tx, booking.TransactionID)

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking detail retrieved", gin.H{
		"booking":     booking,
		"transaction": tx,
	}))
}