package controllers

import (
	"net/http"
	"strconv"

	"app/models"
	"app/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(s *services.UserService) *UserController {
	return &UserController{Service: s}
}

// GetAll godoc
// @Summary      Get All Users
// @Description  Mengambil seluruh daftar data pengguna
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /users [get]
func (u *UserController) GetAll(c *gin.Context) {
	users, err := u.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetByID godoc
// @Summary      Get User Detail
// @Description  Mengambil data detail pengguna berdasarkan ID
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /users/{id} [get]
func (u *UserController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
		return
	}

	user, err := u.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Create godoc
// @Summary      Create User
// @Description  Menambahkan pengguna baru ke dalam database
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.UserCreateRequest true "Data User Baru"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /users [post]
func (u *UserController) Create(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser := c.GetString("username")
	user, err := u.Service.Create(req, authUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User berhasil dibuat", "data": user})
}

// Update godoc
// @Summary      Update User (Partial)
// @Description  Memperbarui data pengguna secara parsial berdasarkan ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Param        request body models.UserUpdateRequest true "Data Perubahan User"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /users/{id} [put]
func (u *UserController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser := c.GetString("username")
	err = u.Service.Update(id, req, authUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User berhasil diperbarui"})
}

// Delete godoc
// @Summary      Delete User
// @Description  Menghapus akun pengguna berdasarkan ID (Auto logout jika menghapus akun sendiri)
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /users/{id} [delete]
func (u *UserController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
		return
	}

	currentUserID := c.GetInt("user_id")

	err = u.Service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if id == currentUserID {
		c.JSON(http.StatusOK, gin.H{
			"message": "Akun Anda sendiri telah berhasil dihapus. Anda telah otomatis logout.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
