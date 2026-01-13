package tests

import (
	"backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllBarang(t *testing.T) {
	r := GetRouter()

	t.Run("Admin can access", func(t *testing.T) {
		// Seed barang
		SeedBarang(t, r)			

		token := LoginAsAdmin(t, r)
		w := PerformRequest(r, "GET", "/api/v1/admin/barang", nil, token)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "nama_barang")
	})

	t.Run("User forbidden", func(t *testing.T) {
		token := LoginAsUser(t, r)
		w := PerformRequest(r, "GET", "/api/v1/admin/barang", nil, token)

		assert.Equal(t, 403, w.Code)
	})

	t.Run("Unauthenticated forbidden", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/v1/admin/barang", nil, "")
		assert.Equal(t, 401, w.Code)
	})
}

func TestAdminGetBarangByID(t *testing.T) {
	r := GetRouter()

	SeedBarang(t, r)
	t.Run("Admin can access", func(t *testing.T) {
		token := LoginAsAdmin(t, r)
		w := PerformRequest(r, "GET", "/api/v1/admin/barang/1", nil, token)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("User forbidden", func(t *testing.T) {
		token := LoginAsUser(t, r)
		w := PerformRequest(r, "GET", "/api/v1/admin/barang/1", nil, token)
		assert.Equal(t, 403, w.Code)
	})

	t.Run("Unauthenticated forbidden", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/v1/admin/barang/1", nil, "")
		assert.Equal(t, 401, w.Code)
	})
}


func TestAdminCreateBarang(t *testing.T) {
	r := GetRouter()

	token := LoginAsAdmin(t, r)

	SeedBarang(t, r)

	payload := map[string]string{
		"nama_barang": "Komputer LCD",
		"stok": "10",
		"image_url": "https://example.com/image.jpg",
	}

	w := PerformRequest(r, "POST", "/api/v1/admin/barang", payload, token)
	assert.Equal(t, 201, w.Code)

	var barang models.Barang
	err := GetDB().
		Where("nama_barang = ?", "Komputer LCD").
		First(&barang).Error

	assert.NoError(t, err)
	assert.Equal(t, "Komputer LCD", barang.Nama_barang)
	assert.Equal(t, "10", barang.Stok)
	assert.Equal(t, "https://example.com/image.jpg", barang.ImageURL)

}

func TestAdminUpdateBarang(t *testing.T){
	r := GetRouter()

	token := LoginAsAdmin(t, r)

	SeedBarang(t, r)

	payload := map[string]string{
		"nama_barang": "Komputer LCD",
		"stok": "10",
		"image_url": "https://example.com/image.jpg",
	}

	w := PerformRequest(r, "PUT", "/api/v1/admin/barang/1", payload, token)
	assert.Equal(t, 200, w.Code)

	var barang models.Barang
	err := GetDB().
		Where("nama_barang = ?", "Komputer LCD").
		First(&barang).Error
	assert.NoError(t, err)
	assert.Equal(t, "Komputer LCD", barang.Nama_barang)
	assert.Equal(t, "10", barang.Stok)
	assert.Equal(t, "https://example.com/image.jpg", barang.ImageURL)
}

