package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"mobile-api/api"
	"mobile-api/config"
	"mobile-api/api/voltras"
	// "mobile-api/db"
)

func main() {
	cfg := config.LoadConfig()

	database := config.ConnectDB(cfg)
	voltras.InitClient(cfg.VoltrasURL, cfg.VoltrasOfficeCode, cfg.VoltrasToken)

	// Auto migrate all tables
	// db.Migrate(database)

	r := gin.Default()

	// CORS middleware (optional - enable if needed by mobile app)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup all routes — pass jwtSecret from config (no os.Getenv anywhere)
	api.SetupRoutes(r, database, cfg.JWTSecret)

	log.Printf("Server running on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}