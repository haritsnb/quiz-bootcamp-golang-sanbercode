package middlewares

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// In-Memory Token Blacklist (Thread-Safe)
var (
	BlacklistedTokens = make(map[string]time.Time)
	blacklistMutex    sync.RWMutex
)

func BlacklistToken(tokenString string, expTime time.Time) {
	// daftarkan token ke blacklist
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()
	BlacklistedTokens[tokenString] = expTime
}

func IsTokenBlacklisted(tokenString string) bool {
	// Cek apakah token sudah di-blacklist
	blacklistMutex.RLock()
	defer blacklistMutex.RUnlock()
	exp, exists := BlacklistedTokens[tokenString]
	if !exists {
		return false
	}
	// Jika token sudah expired secara alami, hapus dari memori blacklist
	if time.Now().After(exp) {
		delete(BlacklistedTokens, tokenString)
		return false
	}
	return true
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header dibutuhkan"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format Authorization harus 'Bearer <token>'"})
			c.Abort()
			return
		}

		if IsTokenBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token sudah tidak valid karena telah logout. Silakan login kembali."})
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau telah kedaluwarsa"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Payload token invalid"})
			c.Abort()
			return
		}

		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Set("token_string", tokenString)
		c.Set("token_exp", int64(claims["exp"].(float64)))
		c.Next()
	}
}
