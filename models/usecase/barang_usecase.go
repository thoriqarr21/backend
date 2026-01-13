package usecase

import (
	"backend/models"
	"backend/models/repository"
	"errors"
)

type BarangUsecase interface {
	GetAllBarang() ([]models.Barang, error)
	GetBarangByID(id uint) (*models.Barang, error)
	CreateBarang(req models.CreateBarangRequest) (*models.Barang, error)
	UpdateBarang(id uint, req models.UpdateBarangRequest) (*models.Barang, error)
	DeleteBarang(id uint) error
}

type barangUsecase struct {
	repo repository.BarangRepository
}

func NewBarangUsecase(repo repository.BarangRepository) BarangUsecase {
	return &barangUsecase{repo: repo}
}

func (b *barangUsecase) GetAllBarang() ([]models.Barang, error) {
	return b.repo.GetAllBarang()
}

func (b *barangUsecase) GetBarangByID(id uint) (*models.Barang, error) {
	return b.repo.FindByID(id)
}

func (b *barangUsecase) CreateBarang(req models.CreateBarangRequest) (*models.Barang, error) {
	// Check if username already exists
	if _, err := b.repo.FindByNama_barang(req.Nama_barang); err == nil {
		return nil, errors.New("username already exists")
	}

	barang := &models.Barang{
		Nama_barang: req.Nama_barang,
		Stok:        req.Stok,
		ImageURL:    req.ImageURL,
	}

	if err := b.repo.Create(barang); err != nil {
		return nil, err
	}

	return barang, nil
}

func (u *barangUsecase) UpdateBarang(id uint, req models.UpdateBarangRequest) (*models.Barang, error) {
	barang, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Nama_barang != nil {
		barang.Nama_barang = *req.Nama_barang
	}
	if req.Stok != nil {
		barang.Stok = *req.Stok
	}
	if req.ImageURL != nil {
		barang.ImageURL = *req.ImageURL
	}

	if err := u.repo.Update(barang); err != nil {
		return nil, err
	}

	return barang, nil
}

func (u *barangUsecase) DeleteBarang(id uint) error {
	_, err := u.repo.FindByID(id)
	if err != nil {
		return errors.New("barang not found")
	}
	return u.repo.Delete(id)
}