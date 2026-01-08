package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	r := GetRouter()

	t.Run("Success login", func(t *testing.T) {
		payload := map[string]string{
			"username": "admin",
			"password": "admin123",
		}

		w := PerformRequest(r, "POST", "/api/v1/auth/login", payload, "")
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Fail wrong password", func(t *testing.T) {
		payload := map[string]string{
			"username": "admin",
			"password": "wrong",
		}

		w := PerformRequest(r, "POST", "/api/v1/auth/login", payload, "")
		assert.Equal(t, 401, w.Code)
	})
}
