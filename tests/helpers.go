package tests

import (
	"backend/config"
	"backend/models"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func seedUsers(r *gin.Engine) {
	if r == nil {
		panic("router cannot be nil when calling seedUsers")
	}
	
	admin := map[string]string{
		"username":  "admin",
		"email":     "admin@ezytix.com",
		"password":  "admin123",
		"full_name": "Admin",
		"role":      "admin",
		"phone":     "0893344443",
	}

	user := map[string]string{
		"username":  "user",
		"email":     "user@ezytix.com",
		"password":  "user123",
		"full_name": "User",
		"role":      "user",
		"phone":     "0893344443",
	}

	PerformRequest(r, "POST", "/api/v1/auth/register", admin, "")
	PerformRequest(r, "POST", "/api/v1/auth/register", user, "")
}

func LoginAsAdmin(t *testing.T, r http.Handler) string {
	return login(t, r, "admin", "admin123")
}

func LoginAsUser(t *testing.T, r http.Handler) string {
	return login(t, r, "user", "user123")
}

func login(t *testing.T, r http.Handler, username, password string) string {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	w := PerformRequest(r, "POST", "/api/v1/auth/login", payload, "")
	assert.Equal(t, 200, w.Code)

	res := ParseResponse(t, w)
	token := res["token"].(string)
	return token
}

func RegisterAndGetUserID(
	t *testing.T,
	r http.Handler,
	username, email string,
) uint {

	payload := map[string]string{
		"username":  username,
		"email":     email,
		"password":  "pass123",
		"full_name": "Temp",
		"role":      "user",
		"phone":     "0893344443",
	}

	w := PerformRequest(r, "POST", "/api/v1/auth/register", payload, "")
	assert.Equal(t, 201, w.Code)

	// Get database connection from config
	cfg := config.LoadConfig()
	db := config.ConnectDB(cfg)
	
	// Query the database
	var user models.User
	err := db.
		Where("username = ?", username).
		First(&user).
		Error

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	return user.ID
}
