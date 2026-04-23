package tests

import (
	"mobile-api/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAdminUpdateUser(t *testing.T) {
	r := GetRouter()

	userID := RegisterAndGetUserID(
		t,
		r,
		"updatetest",
		"updatetest@test.com",
	)

	token := LoginAsAdmin(t, r)

	payload := map[string]string{
		"email":     "updated@test.com",
		"full_name": "Updated",
	}

	w := PerformRequest(
		r,
		"PUT",
		fmt.Sprintf("/api/v1/admin/users/%d", userID),
		payload,
		token,
	)

	assert.Equal(t, 200, w.Code)

	// 🔥 VERIFIKASI LANGSUNG KE DATABASE
	var user models.User
	err := GetDB().First(&user, userID).Error

	assert.NoError(t, err)
	assert.Equal(t, "updatetest", user.Username)
	assert.Equal(t, "updated@test.com", user.Email)
	assert.Equal(t, "Updated", user.FullName)
	// assert.Equal(t, "user", user.Role)
	assert.Equal(t, "0893344443", user.Phone)
}

func TestAdminCreateUser(t *testing.T) {
	r := GetRouter()

	token := LoginAsAdmin(t, r)

	payload := map[string]string{
		"username":  "createduser",
		"email":     "created@test.com",
		"full_name": "Created",
		"password":  "password123",
		"phone":     "0893344443",
	}

	w := PerformRequest(r, "POST", "/api/v1/admin/users", payload, token)
	assert.Equal(t, 201, w.Code)

	var user models.User
	err := GetDB().
		Where("email = ?", "created@test.com").
		First(&user).Error

	assert.NoError(t, err)
	assert.Equal(t, "createduser", user.Username)
	assert.Equal(t, "Created", user.FullName)
	assert.Equal(t, models.RoleUser, user.Role)
	assert.Equal(t, "0893344443", user.Phone)
}

func TestAdminDeleteUser(t *testing.T) {
	r := GetRouter()

	// Buat user dulu
	userID := RegisterAndGetUserID(
		t,
		r,
		"deletetest",
		"deletetest@test.com",
	)

	token := LoginAsAdmin(t, r)

	// Delete user
	w := PerformRequest(
		r,
		"DELETE",
		fmt.Sprintf("/api/v1/admin/users/%d", userID),
		nil,
		token,
	)

	assert.Equal(t, 200, w.Code)

	// Pastikan user benar-benar terhapus
	var user models.User
	err := GetDB().First(&user, userID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestAdminResetPassword(t *testing.T) {
	r := GetRouter()

	userID := RegisterAndGetUserID(
		t,
		r,
		"resettest",
		"resettest@test.com",
	)

	token := LoginAsAdmin(t, r)

	payload := map[string]string{
		"new_password": "newpassword123",
	}

	w := PerformRequest(
		r,
		"POST",
		fmt.Sprintf("/api/v1/admin/users/%d/reset-password", userID),
		payload,
		token,
	)

	assert.Equal(t, 200, w.Code)

	// Verifikasi password berubah
	loginPayload := map[string]string{
		"username": "resettest",
		"password": "newpassword123",
	}

	w = PerformRequest(r, "POST", "/api/v1/auth/login", loginPayload, "")
	assert.Equal(t, 200, w.Code)
}