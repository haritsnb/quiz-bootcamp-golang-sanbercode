package controllers

import (
	"net/http"
	"strconv"

	"app/models"
	"app/services"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	Service *services.CategoryService
}

func NewCategoryController(s *services.CategoryService) *CategoryController {
	return &CategoryController{Service: s}
}

// GetAll godoc
// @Summary      Get All Categories
// @Description  Mengambil seluruh daftar kategori buku
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /categories [get]
func (cc *CategoryController) GetAll(c *gin.Context) {
	data, err := cc.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetByID godoc
// @Summary      Get Category Detail
// @Description  Mengambil detail kategori berdasarkan ID
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /categories/{id} [get]
func (cc *CategoryController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID category tidak valid"})
		return
	}

	cat, err := cc.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cat})
}

// Create godoc
// @Summary      Create Category
// @Description  Menambahkan kategori baru
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CategoryCreateRequest true "Data Kategori Baru"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /categories [post]
func (cc *CategoryController) Create(c *gin.Context) {
	var req models.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser := c.GetString("username")
	cat, err := cc.Service.Create(req, authUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Category berhasil dibuat", "data": cat})
}

// Update godoc
// @Summary      Update Category
// @Description  Memperbarui data kategori berdasarkan ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Param        request body models.CategoryUpdateRequest true "Data Perubahan Kategori"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /categories/{id} [put]
func (cc *CategoryController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID category tidak valid"})
		return
	}

	var req models.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser := c.GetString("username")
	err = cc.Service.Update(id, req, authUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category berhasil diperbarui"})
}

// Delete godoc
// @Summary      Delete Category
// @Description  Menghapus data kategori berdasarkan ID
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /categories/{id} [delete]
func (cc *CategoryController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID category tidak valid"})
		return
	}

	err = cc.Service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category berhasil dihapus"})
}

// GetBooks godoc
// @Summary      Get Books by Category ID
// @Description  Mengambil seluruh buku yang berada dalam kategori tertentu
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /categories/{id}/books [get]
func (cc *CategoryController) GetBooks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID category tidak valid"})
		return
	}

	books, err := cc.Service.GetBooksByCategoryID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Format image_url menjadi Full URL
	for i := range books {
		books[i].ImageURL = formatImageURL(c, books[i].ImageURL)
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}
