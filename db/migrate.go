package db

import (
	"log"

	"gorm.io/gorm"

	"mobile-api/models"
)

func Migrate(db *gorm.DB) {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		// Master data
		&models.Country{},
		&models.Airport{},
		&models.Airline{},
		&models.PaymentMethod{},
		&models.VABank{},

		// Contact
		&models.Contact{},

		// Flight
		&models.FlightTransaction{},
		&models.FlightBooking{},
		&models.FlightPassenger{},
		&models.FlightSegment{},
		&models.FlightRemark{},

		// Balance
		&models.PhoneBalanceTransaction{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed successfully")
}