package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	r := GetRouter()

	t.Run("Admin can access", func(t *testing.T) {
		token := LoginAsAdmin(t, r)
		w := PerformRequest(r, "GET", "/api/v1/admin/users", nil, token)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("User forbidden", func(t *testing.T) {
		token := LoginAsUser(t, r)
		w := PerformRequest(r, "GET", "/api/v1/admin/users", nil, token)
		assert.Equal(t, 403, w.Code)
	})
}

func TestGetUserByID(t *testing.T) {
	r := GetRouter()

	t.Run("Admin can access", func(t *testing.T) {
		token := LoginAsAdmin(t, r)
		w := PerformRequest(r, "GET", "/api/v1/admin/users/1", nil, token)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("User forbidden", func(t *testing.T) {
		token := LoginAsUser(t, r)
		w := PerformRequest(r, "GET", "/api/v1/admin/users/1", nil, token)
		assert.Equal(t, 403, w.Code)
	})
}

func TestGetUserProfile(t *testing.T) {
	r := GetRouter()

	t.Run("User can access own profile", func(t *testing.T) {
		token := LoginAsUser(t, r)
		w := PerformRequest(r, "GET", "/api/v1/users/profile", nil, token)
		assert.Equal(t, 200, w.Code)
	})
}

func TestUpdateProfile(t *testing.T) {
	r := GetRouter()

	t.Run("User can update own profile", func(t *testing.T) {

		username := "testuser"
		email := "testuser@example.com"
		userID := RegisterAndGetUserID(t, r, username, email)

		token := login(t, r, username, "pass123")

		payload := map[string]string{
			"full_name": "Updated Name",
			"email":     "updated@gmail.com",
			"phone":     "08123456789",
		}

		w := PerformRequest(
			r,
			"PUT",
			fmt.Sprintf("/api/v1/users/profile/%d", userID),
			payload,
			token,
		)

		assert.Equal(t, 200, w.Code)
	})
}

func TestChangePassword(t *testing.T) {
	r := GetRouter()

	t.Run("User can change password", func(t *testing.T) {

		username := "testuser"
		email := "testuser@example.com"
		userID := RegisterAndGetUserID(t, r, username, email)
		_ = userID 

		token := login(t, r, username, "pass123")

		payload := map[string]string{
			"old_password": "pass123",
			"new_password": "newpass123",
		}

		w := PerformRequest(
			r,
			"POST",
			"/api/v1/users/change-password",
			payload,
			token,
		)

		assert.Equal(t, 200, w.Code)
	})
}

