package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"app/models"
	"app/services"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	Service *services.BookService
}

func NewBookController(s *services.BookService) *BookController {
	return &BookController{Service: s}
}

// Helper untuk menyusun URL dinamis sesuai protokol dan host client/server
func formatImageURL(c *gin.Context, filename *string) *string {
	if filename == nil || *filename == "" {
		return nil
	}

	// Jika sudah berbentuk URL eksternal (misal: https://...), biarkan apa adanya
	if strings.HasPrefix(*filename, "http://") || strings.HasPrefix(*filename, "https://") {
		return filename
	}

	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	fullURL := fmt.Sprintf("%s://%s/storages/%s", scheme, c.Request.Host, *filename)
	return &fullURL
}

func (bc *BookController) GetAll(c *gin.Context) {
	data, err := bc.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format image_url menjadi Full URL pada setiap buku
	for i := range data {
		data[i].ImageURL = formatImageURL(c, data[i].ImageURL)
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (bc *BookController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID book tidak valid"})
		return
	}

	book, err := bc.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	book.ImageURL = formatImageURL(c, book.ImageURL)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (bc *BookController) Create(c *gin.Context) {
	var req models.BookCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validasi form gagal: " + err.Error()})
		return
	}

	file, err := c.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file gambar: " + err.Error()})
		return
	}

	authUser := c.GetString("username")

	book, err := bc.Service.Create(req, file, authUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan file fisik ke folder storages/
	if file != nil && book.ImageURL != nil {
		diskPath := filepath.Join("storages", *book.ImageURL)
		if err := c.SaveUploadedFile(file, diskPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file fisik: " + err.Error()})
			return
		}
	}

	book.ImageURL = formatImageURL(c, book.ImageURL)
	c.JSON(http.StatusCreated, gin.H{"message": "Book berhasil dibuat", "data": book})
}

func (bc *BookController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID book tidak valid"})
		return
	}

	var req models.BookUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validasi form gagal: " + err.Error()})
		return
	}

	file, err := c.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file gambar: " + err.Error()})
		return
	}

	authUser := c.GetString("username")

	book, err := bc.Service.Update(id, req, file, authUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Simpan file fisik baru ke folder storages/
	if file != nil && book.ImageURL != nil {
		diskPath := filepath.Join("storages", *book.ImageURL)
		if err := c.SaveUploadedFile(file, diskPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file fisik baru: " + err.Error()})
			return
		}
	}

	book.ImageURL = formatImageURL(c, book.ImageURL)
	c.JSON(http.StatusOK, gin.H{"message": "Book berhasil diperbarui", "data": book})
}

func (bc *BookController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID book tidak valid"})
		return
	}

	err = bc.Service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book beserta file image berhasil dihapus"})
}
