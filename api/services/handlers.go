package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mobile-api/models"
	"mobile-api/utils"
)

// GET /api/v1/services/countries
// func GetCountriesHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var countries []models.Country
// 		if err := db.Where("is_active = ?", true).Order("name ASC").Find(&countries).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch countries", err.Error()))
// 			return
// 		}
// 		c.JSON(http.StatusOK, utils.SuccessResponse("Countries retrieved", gin.H{
// 			"countries": countries,
// 			"total":     len(countries),
// 		}))
// 	}
// }

func GetCountriesHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var countries []models.Country
        
        // Ambil SEMUA data tanpa filter is_active
        if err := db.Order("name ASC").Find(&countries).Error; err != nil {
            c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch countries", err.Error()))
            return
        }

        c.JSON(http.StatusOK, utils.SuccessResponse("Countries retrieved", gin.H{
            "countries": countries,
            "total":     len(countries),
        }))
    }
}

// GET /api/v1/services/airports?search=CGK&country_id=1
func GetAirportsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		search    := c.Query("search")
		countryID := c.Query("country_id")

		query := db.Preload("Country").Where("is_active = ?", true).Order("name ASC")
		if search != "" {
			like := "%" + search + "%"
			query = query.Where("iata_code ILIKE ? OR name ILIKE ? OR city ILIKE ?", like, like, like)
		}
		if countryID != "" {
			query = query.Where("country_id = ?", countryID)
		}

		var airports []models.Airport
		if err := query.Find(&airports).Error; err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch airports", err.Error()))
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse("Airports retrieved", gin.H{
			"airports": airports,
			"total":    len(airports),
		}))
	}
}

// GET /api/v1/services/airlines
func GetAirlinesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var airlines []models.Airline
		if err := db.Order("name ASC").Find(&airlines).Error; err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch airlines", err.Error()))
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse("Airlines retrieved", gin.H{
			"airlines": airlines,
			"total":    len(airlines),
		}))
	}
}

// GET /api/v1/services/airlines/:id/logo
func GetAirlineLogoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", ""))
			return
		}
		var airline models.Airline
		if err := db.First(&airline, id).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Airline not found", ""))
			return
		}
		if airline.LogoBase64 == "" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Logo not available", ""))
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse("Airline logo retrieved", gin.H{
			"id":          airline.ID,
			"name":        airline.Name,
			"iata_code":   airline.IATACode,
			"logo_base64": airline.LogoBase64,
		}))
	}
}

// GET /api/v1/services/va-banks
func GetVABanksHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var banks []models.VABank
		if err := db.Preload("PaymentMethod").
			Where("is_active = ?", true).
			Order("bank_name ASC").
			Find(&banks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch VA banks", err.Error()))
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse("VA banks retrieved", gin.H{
			"banks": banks,
			"total": len(banks),
		}))
	}
}

// GET /api/v1/services/va-banks/:id/logo
func GetVABankLogoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", ""))
			return
		}
		var bank models.VABank
		if err := db.First(&bank, id).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("VA bank not found", ""))
			return
		}
		if bank.LogoBase64 == "" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Logo not available", ""))
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse("VA bank logo retrieved", gin.H{
			"id":          bank.ID,
			"bank_name":   bank.BankName,
			"bank_code":   bank.BankCode,
			"logo_base64": bank.LogoBase64,
		}))
	}
}

// POST /api/v1/services/payment/transactions
func CreateTransactionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateTransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}

		var booking models.FlightBooking
		if err := db.First(&booking, req.BookingID).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
			return
		}

		var pm models.PaymentMethod
		if err := db.First(&pm, req.PaymentMethodID).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Payment method not found", ""))
			return
		}

		vaNumber := ""
		if req.BankCode != "" {
			var bank models.VABank
			if err := db.Where("bank_code = ? AND is_active = ?", req.BankCode, true).First(&bank).Error; err == nil {
				vaNumber = fmt.Sprintf("%s%010d", req.BankCode, booking.ID)
			}
		}

		c.JSON(http.StatusCreated, utils.SuccessResponse("Transaction created", gin.H{
			"booking_id":        booking.ID,
			"payment_method_id": pm.ID,
			"payment_method":    pm.Name,
			"bank_code":         req.BankCode,
			"va_number":         vaNumber,
			"amount":            req.Amount,
			"admin_fee":         req.AdminFee,
			"markup":            req.Markup,
			"status":            "pending",
		}))
	}
}

// POST /api/v1/services/transactions/check
func CheckTransactionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CheckTransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse("Transaction status retrieved", gin.H{
			"external_id": req.ExternalID,
			"status":      "pending",
		}))
	}
}

// POST /api/v1/services/documents/send
func SendBookingDocumentsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SendDocumentsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}
		var booking models.FlightBooking
		if err := db.Preload("Passengers").Preload("Segments").First(&booking, req.BookingID).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse("Documents sent successfully", gin.H{
			"booking_id": booking.ID,
			"pnr":        booking.PNRID,
			"sent_to":    req.Email,
		}))
	}
}

// POST /api/v1/services/documents/print
func PrintDocumentsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PrintDocumentsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
			return
		}
		var booking models.FlightBooking
		if err := db.Preload("Passengers").Preload("Segments").First(&booking, req.BookingID).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found", ""))
			return
		}
		docType := req.DocType
		if docType == "" {
			docType = "ITINERARY"
		}
		c.JSON(http.StatusOK, utils.SuccessResponse("Document ready", gin.H{
			"booking_id": booking.ID,
			"pnr":        booking.PNRID,
			"doc_type":   docType,
			"doc_url":    fmt.Sprintf("/uploads/docs/%s_%s.pdf", docType, booking.PNRID),
		}))
	}
}

// GET /api/v1/services/balance/:phone
func GetBalanceHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		phone := c.Param("phone")

		var txs []models.PhoneBalanceTransaction
		db.Where("phone = ?", phone).Order("created_at DESC").Limit(10).Find(&txs)

		var balance float64
		for _, tx := range txs {
			if tx.Type == "TOPUP" {
				balance += tx.Amount
			} else if tx.Type == "DEDUCT" {
				balance -= tx.Amount
			}
		}

		c.JSON(http.StatusOK, utils.SuccessResponse("Balance retrieved", gin.H{
			"phone":        phone,
			"balance":      balance,
			"currency":     "IDR",
			"transactions": txs,
		}))
	}
}