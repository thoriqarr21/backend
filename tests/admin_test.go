package tests

import (
	"backend/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
