package api

import (
	"backend/controllers"
	"backend/controllers/middleware"
	"backend/models"
	"backend/models/repository"
	"backend/models/usecase"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}
	r.Static("/uploads", "./uploads")

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userController := controllers.NewUserController(userUsecase)

	barangRepo := repository.NewBarangRepository(db)
	barangUsecase := usecase.NewBarangUsecase(barangRepo)
	barangController := controllers.NewBarangController(barangUsecase)

	historyRepo := repository.NewTVMReportHistoryRepository(db)
	historyUsecase := usecase.NewTVMReportHistoryUsecase(historyRepo)
	historyController := controllers.NewTVMReportHistoryController(historyUsecase)

	tvmRepo := repository.NewTVMReportRepository(db)
	tvmUsecase := usecase.NewTVMReportUsecase(tvmRepo, barangRepo, historyRepo)
	tvmController := controllers.NewTVMReportController(tvmUsecase)


	// --- API v1 GROUP ---
	v1 := r.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// Protected auth routes
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

		barang := v1.Group("/barang")
		barang.Use(middleware.AuthMiddleware())
		{
			barang.GET("/", barangController.GetAllBarang)
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
			admin.GET("/dashboard", tvmController.GetDashboard)
			admin.GET("/reports", tvmController.GetAllReports)
			admin.PUT("/reports/:id", tvmController.UpdateReport)
			admin.PATCH("/reports/:id/open", tvmController.OpenReport)
			admin.DELETE("/reports/:id", tvmController.DeleteReport)
			admin.GET("/barang", barangController.GetAllBarang)
			admin.GET("/barang/:id", barangController.GetBarangByID)
			admin.POST("/barang", barangController.CreateBarang)
			admin.PUT("/barang/:id", barangController.UpdateBarang)
			admin.DELETE("/barang/:id", barangController.DeleteBarang)
		}
		// routes/api.go
		teknisi := v1.Group("/teknisi")
		teknisi.Use(middleware.AuthMiddleware())
		teknisi.Use(middleware.RoleMiddleware(db, models.RoleTeknisi))
		{
			teknisi.GET("/reports/open", tvmController.GetOpenReports)
			teknisi.PATCH("/reports/:id/take", tvmController.TakeReport)
			teknisi.GET("/reports/my", tvmController.GetMyReportsAsTechnician)
			teknisi.PUT("/reports/:id/resolve", tvmController.ResolveReport)
		}

		// Protected routes - TVM Reports (User & Admin)
		reports := v1.Group("/reports")
		reports.Use(middleware.AuthMiddleware())
		{
			// Endpoint ini sekarang mengharapkan multipart/form-data untuk upload file
			reports.POST("/", tvmController.CreateReport) 
			reports.POST("", tvmController.CreateReport) 
			reports.GET("/:id/history", historyController.GetByReportID)
			reports.GET("/history", historyController.GetAll)
			reports.GET("/:id/historyuser", historyController.FindByChangedBy)
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