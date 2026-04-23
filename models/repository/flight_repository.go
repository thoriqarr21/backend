package repository

import (
	"mobile-api/models"

	"gorm.io/gorm"
)

type FlightRepository interface {
	CreateContact(contact *models.Contact) error
	FindOrCreateContact(email string, contact *models.Contact) error
	CreateTransaction(tx *models.FlightTransaction) error
	CreateBooking(booking *models.FlightBooking) error
	FindBookingByID(id uint) (*models.FlightBooking, error)
	FindBookingByPNR(pnr string) (*models.FlightBooking, error)
	UpdateBookingStatus(id uint, status string) error
	CreatePassenger(pax *models.FlightPassenger) error
	CreateSegment(seg *models.FlightSegment) error
	CreateRemark(remark *models.FlightRemark) error
}

type flightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) CreateContact(contact *models.Contact) error {
	return r.db.Create(contact).Error
}

func (r *flightRepository) FindOrCreateContact(email string, contact *models.Contact) error {
	return r.db.Where("email = ?", email).FirstOrCreate(contact).Error
}

func (r *flightRepository) CreateTransaction(tx *models.FlightTransaction) error {
	return r.db.Create(tx).Error
}

func (r *flightRepository) CreateBooking(booking *models.FlightBooking) error {
	return r.db.Create(booking).Error
}

func (r *flightRepository) FindBookingByID(id uint) (*models.FlightBooking, error) {
	var booking models.FlightBooking
	err := r.db.
		Preload("Passengers").
		Preload("Segments").
		Preload("Remarks").
		First(&booking, id).Error
	return &booking, err
}

func (r *flightRepository) FindBookingByPNR(pnr string) (*models.FlightBooking, error) {
	var booking models.FlightBooking
	err := r.db.
		Preload("Passengers").
		Preload("Segments").
		Where("pnrid = ?", pnr).
		First(&booking).Error
	return &booking, err
}

func (r *flightRepository) UpdateBookingStatus(id uint, status string) error {
	return r.db.Model(&models.FlightBooking{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *flightRepository) CreatePassenger(pax *models.FlightPassenger) error {
	return r.db.Create(pax).Error
}

func (r *flightRepository) CreateSegment(seg *models.FlightSegment) error {
	return r.db.Create(seg).Error
}

func (r *flightRepository) CreateRemark(remark *models.FlightRemark) error {
	return r.db.Create(remark).Error
}