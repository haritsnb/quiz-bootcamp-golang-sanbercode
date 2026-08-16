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

func formatImageURL(c *gin.Context, filename *string) *string {
	if filename == nil || *filename == "" {
		return nil
	}

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

// GetAll godoc
// @Summary      Get All Books
// @Description  Mengambil seluruh data buku beserta Full URL cover image
// @Tags         Books
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /books [get]
func (bc *BookController) GetAll(c *gin.Context) {
	data, err := bc.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range data {
		data[i].ImageURL = formatImageURL(c, data[i].ImageURL)
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetByID godoc
// @Summary      Get Book Detail
// @Description  Mengambil detail data buku berdasarkan ID
// @Tags         Books
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Book ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /books/{id} [get]
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

// Create godoc
// @Summary      Create Book
// @Description  Menambahkan buku baru (Multipart Form) dengan upload file cover dan auto thickness
// @Tags         Books
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        title formData string true "Judul Buku"
// @Param        description formData string false "Deskripsi"
// @Param        release_year formData int true "Tahun Rilis (1980 - 2024)"
// @Param        price formData int true "Harga Buku"
// @Param        total_page formData int true "Total Halaman (Konversi ketebalan otomatis)"
// @Param        category_id formData int true "Category ID"
// @Param        created_by formData string false "Dibuat Oleh (Opsional)"
// @Param        modified_by formData string false "Diupdate Oleh (Opsional)"
// @Param        image formData file false "File Gambar Cover (.jpg/.png)"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /books [post]
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

// Update godoc
// @Summary      Update Book (Partial)
// @Description  Memperbarui data buku atau mengganti file cover
// @Tags         Books
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Book ID"
// @Param        title formData string false "Judul Buku"
// @Param        description formData string false "Deskripsi"
// @Param        release_year formData int false "Tahun Rilis (1980 - 2024)"
// @Param        price formData int false "Harga"
// @Param        total_page formData int false "Total Halaman"
// @Param        category_id formData int false "Category ID"
// @Param        modified_by formData string false "Diupdate Oleh"
// @Param        image formData file false "File Gambar Baru"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /books/{id} [put]
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

// Delete godoc
// @Summary      Delete Book
// @Description  Menghapus data buku dan menghapus file gambar fisiknya dari storage
// @Tags         Books
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Book ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /books/{id} [delete]
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
