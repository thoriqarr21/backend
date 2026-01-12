package tests

import (
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