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

	// Initialize Barang layers
	barangRepo := repository.NewBarangRepository(db)
	barangUsecase := usecase.NewBarangUsecase(barangRepo)
	barangController := controllers.NewBarangController(barangUsecase)

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Public routes (tanpa auth)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// Protected auth routes (perlu token)
		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware())
		{
			authProtected.POST("/logout", userController.Logout)
		}

		// Protected routes - All authenticated users
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", userController.GetProfile)
			users.PUT("/profile/:id", userController.UpdateUser)
			users.POST("/change-password", userController.ChangePassword)
		}

		// Protected routes - Admin only
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.RoleMiddleware(db, models.RoleAdmin))
		{
			admin.GET("/users", userController.GetAllUsers)
			admin.POST("/users", userController.CreateUser)
			admin.GET("/users/:id", userController.GetUserByID)
			admin.PUT("/users/:id", userController.UpdateUser)
			admin.DELETE("/users/:id", userController.DeleteUser)
			admin.GET("/reports/statistics", tvmController.GetStatistics)
			admin.GET("/reports", tvmController.GetAllReports)
			admin.PUT("/reports/:id", tvmController.UpdateReport)
			admin.DELETE("/reports/:id", tvmController.DeleteReport)
			admin.GET("/barang", barangController.GetAllBarang)
			admin.POST("/barang", barangController.CreateBarang)
			admin.PUT("/barang/:id", barangController.UpdateBarang)
			admin.DELETE("/barang/:id", barangController.DeleteBarang)
		}

		// Protected routes - TVM Reports (User & Admin)
		reports := v1.Group("/reports")
		reports.Use(middleware.AuthMiddleware())
		{
			reports.POST("/", tvmController.CreateReport)
			reports.GET("/", tvmController.GetAllReports)
			reports.PATCH("/:id/status", tvmController.UpdateStatusByPetugas)
			reports.GET("/my", tvmController.GetMyReports)
			reports.GET("/:id", tvmController.GetReportByID)
			reports.PUT("/:id", tvmController.UpdateReport)
			reports.DELETE("/:id", tvmController.DeleteReport)
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