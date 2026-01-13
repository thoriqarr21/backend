package repository

import (
	"backend/models"

	"gorm.io/gorm"
)

type BarangRepository interface {
	GetAllBarang() ([]models.Barang, error)
	FindByID(id uint) (*models.Barang, error)
	FindByNama_barang(nama_barang string) (*models.Barang, error)
	FindByTVMCode(tvmCode string) (*models.Barang, error)
	Create(barang *models.Barang) error
	Update(barang *models.Barang) error
	Delete(id uint) error
	CountBarang() (int64, error)
}

func (r *barangRepository) FindByNama_barang(nama_barang string) (*models.Barang, error) {
	var barang models.Barang
	err := r.db.Where("nama_barang = ?", nama_barang).First(&barang).Error
	if err != nil {
		return nil, err
	}
	return &barang, nil
}

type barangRepository struct {
	db *gorm.DB
}

func NewBarangRepository(db *gorm.DB) BarangRepository {
	return &barangRepository{db: db}
}

func (r *barangRepository) GetAllBarang() ([]models.Barang, error) {
	var barang []models.Barang
	err := r.db.Find(&barang).Error
	return barang, err
}

func (r *barangRepository) GetBarangByID(id uint) (*models.Barang, error) {
	var barang models.Barang
	err := r.db.First(&barang, id).Error
	return &barang, err
}

func (r *barangRepository) Create(barang *models.Barang) error {
	return r.db.Create(barang).Error
}

func (r *barangRepository) Update(barang *models.Barang) error {
	return r.db.Save(barang).Error
}

func (r *barangRepository) Delete(id uint) error {
	return r.db.Delete(&models.Barang{}, id).Error
}

func (r *barangRepository) FindByTVMCode(tvmCode string) (*models.Barang, error) {
	var barang models.Barang
	err := r.db.Where("tvm_code = ?", tvmCode).First(&barang).Error
	if err != nil {
		return nil, err
	}
	return &barang, nil
}

func (r *barangRepository) FindByID(id uint) (*models.Barang, error) {
	var barang models.Barang
	err := r.db.First(&barang, id).Error
	if err != nil {
		return nil, err
	}
	return &barang, nil
}

func (r *barangRepository) CountBarang() (int64, error) {
	var total int64
	err := r.db.Model(&models.Barang{}).Count(&total).Error
	return total, err
}
