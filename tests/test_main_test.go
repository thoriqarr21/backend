package tests

import (
	"backend/api"
	"backend/config"

	// "backend/models"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	router *gin.Engine
	db     *gorm.DB
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	setupTestDB()
	code := m.Run()
	cleanupTestDB()

	os.Exit(code)
}

func setupTestDB() {
	cfg := &config.Config{
		DBHost:     getEnv("TEST_DB_HOST", "localhost"),
		DBPort:     getEnv("TEST_DB_PORT", "5432"),
		DBUser:     getEnv("TEST_DB_USER", "postgres"),
		DBPassword: getEnv("TEST_DB_PASSWORD", "12345"),
		DBName:     "ezytix_db", // DB ASLI
		ServerPort: "8080",
		JWTSecret:  "test-secret",
	}

	os.Setenv("JWT_SECRET", cfg.JWTSecret)

	db = config.ConnectDB(cfg)

	// ❌ JANGAN AutoMigrate

	// Reset data (AMAN)
	db.Exec(`
		TRUNCATE users, tvm_reports, barang
		RESTART IDENTITY CASCADE
	`)

	router = gin.New()
	api.SetupRoutes(router, db)

	seedUsers(router)
}


func cleanupTestDB() {
	db.Exec(`
		TRUNCATE users, tvm_reports, barang
		RESTART IDENTITY CASCADE
	`)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func GetRouter() *gin.Engine {
	if router == nil {
		panic("GetRouter() called but router is nil — TestMain not executed")
	}
	return router
}


func GetDB() *gorm.DB {
	return db
}
