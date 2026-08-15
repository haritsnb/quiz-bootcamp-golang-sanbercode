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

func (cc *CategoryController) GetAll(c *gin.Context) {
	data, err := cc.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

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
