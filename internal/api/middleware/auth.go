package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/affiliate-backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Context keys
const (
	UserIDKey    = "userID"
	UserEmailKey = "userEmail"
	UserRoleKey  = "userRole"
)

// AuthClaims represents the JWT claims from Supabase
type AuthClaims struct {
	jwt.RegisteredClaims
	SessionID string `json:"session_id"`
	// Add other claims Supabase might include if needed
}

// AuthMiddleware validates Supabase JWTs
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		tokenString := parts[1]

		claims := &AuthClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.SupabaseJWTSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			return
		}

		userID := claims.Subject // 'sub' claim is the UserID from Supabase
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		// Set user ID in context for downstream handlers
		c.Set(UserIDKey, userID)

		c.Next()
	}
}