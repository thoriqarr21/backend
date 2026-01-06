package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims contains custom JWT claims
type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware verifies the JWT token and sets the user ID in the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			return
		}

		// Get JWT secret from environment variable
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Check token expiration
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			return
		}

		// Set user ID in context
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}