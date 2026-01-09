// package config

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// type Config struct {
// 	DBHost     string
// 	DBPort     string
// 	DBUser     string
// 	DBPassword string
// 	DBName     string
// 	ServerPort string
// 	JWTSecret  string
// }

// func LoadConfig() *Config {
// 	godotenv.Load()

// 	return &Config{
// 		DBHost:     getEnv("DB_HOST", "172.16.10.82"),
// 		DBPort:     getEnv("DB_PORT", "3306"),
// 		DBUser:     getEnv("DB_USER", "root"),
// 		DBPassword: getEnv("DB_PASSWORD", ""),
// 		DBName:     getEnv("DB_NAME", "mydb"),
// 		ServerPort: getEnv("SERVER_PORT", "8080"),
// 		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-this"),
// 	}
// }

// func getEnv(key, defaultValue string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}
// 	return defaultValue
// }

// func ConnectDB(cfg *Config) *gorm.DB {
// 	dsn := fmt.Sprintf(
// 		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
// 	)

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	log.Println("Database connected successfully")
// 	return db
// }

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	JWTSecret  string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "172.16.10.82"),
		DBPort:     getEnv("DB_PORT", "5432"), // Port PostgreSQL
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "12345"),
		DBName:     getEnv("DB_NAME", "ezytix_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-this"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ConnectDB(cfg *Config) *gorm.DB {
	// DSN format untuk PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
	return db
}