package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	servicesPkg "mobile-api/api/services"
	voltrasPkg "mobile-api/api/voltras"
	"mobile-api/controllers"
	"mobile-api/controllers/middleware"
	"mobile-api/models"
	"mobile-api/models/repository"
	"mobile-api/models/usecase"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, jwtSecret string) {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}
	r.Static("/uploads", "./uploads")

	// ─── CONTROLLERS ─────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)
	userController := controllers.NewUserController(userUsecase)
	bookingController := controllers.NewBookingController(db)

	// ─── API v1 ──────────────────────────────────────────────────────────────
	v1 := r.Group("/api/v1")

	// Public auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	// Auth middleware shorthand
	authMW := middleware.AuthMiddleware(jwtSecret)
	adminMW := middleware.RoleMiddleware(db, models.RoleAdmin)

	// Protected auth routes
	authProtected := v1.Group("/auth")
	authProtected.Use(authMW)
	{
		authProtected.POST("/logout", userController.Logout)
	}

	// User profile routes
	users := v1.Group("/users")
	users.Use(authMW)
	{
		users.GET("/profile", userController.GetProfile)
		users.PUT("/profile/:id", userController.UpdateUser)
		users.POST("/change-password", userController.ChangePassword)
		users.POST("/reset-password", userController.ResetPassword)
	}

	// User booking history
	bookings := v1.Group("/bookings")
	bookings.Use(authMW)
	{
		bookings.GET("", bookingController.GetMyBookings)
		bookings.GET("/:id", bookingController.GetBookingDetail)
	}

	// ─── VOLTRAS FLIGHT ROUTES ───────────────────────────────────────────────
	voltras := v1.Group("/voltras")
	{
		// Public - no auth needed
		voltras.GET("/office-info", voltrasPkg.GetOfficeInformationHandler)

		// Protected
		voltrasAuth := voltras.Group("")
		voltrasAuth.Use(authMW)
		{
			voltrasAuth.POST("/flight/search", voltrasPkg.FlightAvailabilityHandler(db))
			voltrasAuth.POST("/flight/fare/retrieve", voltrasPkg.RetrieveFareHandler)
			voltrasAuth.POST("/flight/book", voltrasPkg.BookFlightHandler(db))
			voltrasAuth.POST("/flight/pnr/retrieve", voltrasPkg.RetrievePNRHandler(db))
			voltrasAuth.POST("/flight/ticket", voltrasPkg.TicketingFlightHandler(db))
			voltrasAuth.POST("/flight/ticket/cancel", voltrasPkg.CancelTicketHandler(db))
			voltrasAuth.POST("/flight/autoticket", voltrasPkg.AutoTicketHandler(db))
			voltrasAuth.GET("/flight/print", voltrasPkg.PrintTicketHandler)
			voltrasAuth.GET("/flight/print-insurance", voltrasPkg.PrintInsuranceHandler)
			voltrasAuth.POST("/flight/advance-retrieve", voltrasPkg.AdvanceRetrieveHandler)
		}
	}

	// ─── SERVICES ROUTES ─────────────────────────────────────────────────────
	svc := v1.Group("/services")
	svc.Use(authMW)
	{
		// Master data
		svc.GET("/countries", servicesPkg.GetCountriesHandler(db))
		svc.GET("/airports", servicesPkg.GetAirportsHandler(db))
		svc.GET("/airlines", servicesPkg.GetAirlinesHandler(db))
		svc.GET("/airlines/:id/logo", servicesPkg.GetAirlineLogoHandler(db))
		svc.GET("/va-banks", servicesPkg.GetVABanksHandler(db))
		svc.GET("/va-banks/:id/logo", servicesPkg.GetVABankLogoHandler(db))

		// Payment
		svc.POST("/payment/transactions", servicesPkg.CreateTransactionHandler(db))
		svc.POST("/transactions/check", servicesPkg.CheckTransactionHandler())

		// Documents
		svc.POST("/documents/send", servicesPkg.SendBookingDocumentsHandler(db))
		svc.POST("/documents/print", servicesPkg.PrintDocumentsHandler(db))

		// Balance
		svc.GET("/balance/:phone", servicesPkg.GetBalanceHandler(db))
	}

	// ─── ADMIN ROUTES ─────────────────────────────────────────────────────────
	admin := v1.Group("/admin")
	admin.Use(authMW, adminMW)
	{
		admin.GET("/users", userController.GetAllUsers)
	}

	// ─── HEALTH CHECK ─────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})
}