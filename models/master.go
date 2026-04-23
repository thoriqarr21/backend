package models

import (
	"time"

	// "gorm.io/gorm"
)

// ─── t_countries ─────────────────────────────────────────────────────────────

type Country struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	ISOCode      string         `gorm:"column:iso_code;uniqueIndex;size:5" json:"iso_code"`
	ISOCode3     string         `gorm:"column:iso_code3;size:5" json:"iso_code3"`
	CurrencyCode string         `gorm:"column:currency_code;size:10" json:"currency_code"`
	PhoneCode    string         `gorm:"column:phone_code;size:10" json:"phone_code"`
	IsActive     bool           `gorm:"column:is_active;default:true" json:"is_active"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (Country) TableName() string { return "t_countries" }

// ─── t_airports ──────────────────────────────────────────────────────────────

type Airport struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:150;not null" json:"name"`
	IATACode  string         `gorm:"column:iata_code;uniqueIndex;size:5" json:"iata_code"`
	City      string         `gorm:"size:100" json:"city"`
	CountryID uint           `gorm:"column:country_id" json:"country_id"`
	Country   *Country       `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	IsActive  bool           `gorm:"column:is_active;default:true" json:"is_active"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (Airport) TableName() string { return "t_airports" }

// ─── t_airlines ──────────────────────────────────────────────────────────────

type Airline struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"size:100;not null" json:"name"`
	IATACode   string `gorm:"column:iata_code;uniqueIndex;size:5" json:"iata_code"`
	LogoBase64 string `gorm:"column:logo_base64;type:text" json:"logo_base64,omitempty"`
}

func (Airline) TableName() string { return "t_airlines" }

// ─── t_payment_methods ───────────────────────────────────────────────────────

type PaymentMethod struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:20;uniqueIndex" json:"code"`
	Name string `gorm:"size:100" json:"name"`
}

func (PaymentMethod) TableName() string { return "t_payment_methods" }

// ─── t_va_banks ──────────────────────────────────────────────────────────────

type VABank struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	PaymentMethodID uint           `gorm:"column:payment_method_id" json:"payment_method_id"`
	PaymentMethod   *PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	BankName        string         `gorm:"column:bank_name;size:100" json:"bank_name"`
	BankCode        string         `gorm:"column:bank_code;size:20;uniqueIndex" json:"bank_code"`
	LogoBase64      string         `gorm:"column:logo_base64;type:text" json:"logo_base64,omitempty"`
	IsActive        bool           `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (VABank) TableName() string { return "t_va_banks" }