// package main

// import (
// 	"ezytix-be-api/config"
// 	"ezytix-be-api/routes"
// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	app := fiber.New()

// 	config.ConnectDB()    // Hubungkan DB
// 	routes.SetupRoutes(app) // Setup Route

// 	app.Listen(":3000")
// }

package main

import (
	"backend/api"
	"backend/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database := config.ConnectDB(cfg)

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	api.SetupRoutes(r, database)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}