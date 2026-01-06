package models

import (
	"time"

	// "gorm.io/gorm"
)

type Barang struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Nama_barang  string      `gorm:"type:varchar(255)" json:"nama_barang"`
	Stok     string          `gorm:"type:varchar(255)" json:"stok"`
	ImageURL string          `gorm:"type:varchar(500);column:image_url" json:"image_url,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ApiResponse struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}


// TableName overrides the default table name
func (Barang) TableName() string {
	return "barang"
}

type CreateBarangRequest struct {
	Nama_barang     string `json:"nama_barang" binding:"required"`
	Stok    		string `json:"stok" binding:"required"`
	ImageURL    	string `json:"image_url"`
}

type UpdateBarangRequest struct {
	Nama_barang     *string `json:"nama_barang"`
	Stok    		*string `json:"stok"`
	ImageURL    	*string `json:"image_url"`
}