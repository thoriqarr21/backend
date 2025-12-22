package api

import (
	"backend/controllers"
	"backend/controllers/middleware"
	"backend/models"
	"backend/models/repository"
	"backend/models/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize user layers
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userController := controllers.NewUserController(userUsecase)

	// Initialize TVM report layers
	tvmRepo := repository.NewTVMReportRepository(db)
	tvmUsecase := usecase.NewTVMReportUsecase(tvmRepo)
	tvmController := controllers.NewTVMReportController(tvmUsecase)

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// Protected routes - All authenticated users
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", userController.GetProfile)
		}

		// Protected routes - Admin only
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.RoleMiddleware(db, models.RoleAdmin))
		{
			admin.GET("/users", userController.GetAllUsers)
			admin.GET("/reports/statistics", tvmController.GetStatistics)
		}

		// Protected routes - TVM Reports (User & Admin)
		reports := v1.Group("/reports")
		reports.Use(middleware.AuthMiddleware())
		{
			reports.POST("/", tvmController.CreateReport)                // User & Admin can create
			reports.GET("/", tvmController.GetAllReports)                // User & Admin can view all
			reports.GET("/my", tvmController.GetMyReports)               // Get my reports only
			reports.GET("/:id", tvmController.GetReportByID)             // Get specific report
			reports.PUT("/:id", tvmController.UpdateReport)              // Update report (own or admin)
			reports.DELETE("/:id", tvmController.DeleteReport)           // Admin only (handled in controller)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Server is running",
		})
	})
}