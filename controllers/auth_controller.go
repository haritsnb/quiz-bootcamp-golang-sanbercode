package controllers

import (
	"net/http"
	"strings"
	"time"

	"app/middlewares"
	"app/models"
	"app/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(s *services.AuthService) *AuthController {
	return &AuthController{AuthService: s}
}

func (a *AuthController) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.AuthService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	expUnix := c.GetInt64("token_exp")
	expTime := time.Unix(expUnix, 0)

	// Simpan token yang sedang aktif ke blacklist
	middlewares.BlacklistToken(tokenString, expTime)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout berhasil. Token telah dibatalkan dan tidak dapat digunakan lagi.",
	})
}
