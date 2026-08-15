package services

import (
	"app/models"
	"app/repositories"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type BookService struct {
	Repo         *repositories.BookRepository
	CategoryRepo *repositories.CategoryRepository
}

func NewBookService(r *repositories.BookRepository, c *repositories.CategoryRepository) *BookService {
	return &BookService{Repo: r, CategoryRepo: c}
}

func createSlug(title string) string {
	s := strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = fmt.Sprintf("book-%d", time.Now().Unix())
	}
	return s
}

func (s *BookService) GetAll() ([]models.Book, error) {
	return s.Repo.FindAll()
}

func (s *BookService) GetByID(id int) (*models.Book, error) {
	return s.Repo.FindByID(id)
}

func (s *BookService) Create(req models.BookCreateRequest, file *multipart.FileHeader, authUser string) (*models.Book, error) {
	if _, err := s.CategoryRepo.FindByID(req.CategoryID); err != nil {
		return nil, fmt.Errorf("category_id %d tidak ditemukan", req.CategoryID)
	}

	thickness := "tipis"
	if req.TotalPage > 100 {
		thickness = "tebal"
	}

	var filenameOnly *string
	if file != nil {
		_ = os.MkdirAll("storages", os.ModePerm)

		ext := filepath.Ext(file.Filename)
		if ext == "" {
			ext = ".jpg"
		}

		slug := createSlug(req.Title)
		filename := fmt.Sprintf("%s-%d%s", slug, time.Now().Unix(), ext)
		filenameOnly = &filename
	}

	creator := "system"
	if req.CreatedBy != nil && *req.CreatedBy != "" {
		creator = *req.CreatedBy
	} else if authUser != "" {
		creator = authUser
	}

	modifier := creator
	if req.ModifiedBy != nil && *req.ModifiedBy != "" {
		modifier = *req.ModifiedBy
	}

	book := &models.Book{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    filenameOnly,
		ReleaseYear: req.ReleaseYear,
		Price:       req.Price,
		TotalPage:   req.TotalPage,
		Thickness:   thickness,
		CategoryID:  req.CategoryID,
		CreatedBy:   &creator,
		ModifiedBy:  &modifier,
	}

	err := s.Repo.Create(book)
	return book, err
}

func (s *BookService) Update(id int, req models.BookUpdateRequest, file *multipart.FileHeader, authUser string) (*models.Book, error) {
	existingBook, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]interface{})

	targetTitle := existingBook.Title
	if req.Title != nil && *req.Title != "" {
		fields["title"] = *req.Title
		targetTitle = *req.Title
	}

	if req.Description != nil {
		fields["description"] = *req.Description
	}

	if req.ReleaseYear != nil {
		if *req.ReleaseYear < 1980 || *req.ReleaseYear > 2024 {
			return nil, errors.New("release_year harus antara 1980 dan 2024")
		}
		fields["release_year"] = *req.ReleaseYear
	}

	if req.Price != nil {
		if *req.Price < 0 {
			return nil, errors.New("price tidak boleh kurang dari 0")
		}
		fields["price"] = *req.Price
	}

	if req.TotalPage != nil {
		if *req.TotalPage < 1 {
			return nil, errors.New("total_page minimal 1")
		}
		fields["total_page"] = *req.TotalPage
		if *req.TotalPage > 100 {
			fields["thickness"] = "tebal"
		} else {
			fields["thickness"] = "tipis"
		}
	}

	if req.CategoryID != nil {
		if _, err := s.CategoryRepo.FindByID(*req.CategoryID); err != nil {
			return nil, fmt.Errorf("category_id %d tidak ditemukan", *req.CategoryID)
		}
		fields["category_id"] = *req.CategoryID
	}

	if file != nil {
		// Hapus file fisik lama dari disk jika ada
		if existingBook.ImageURL != nil && *existingBook.ImageURL != "" {
			oldPath := filepath.Join("storages", *existingBook.ImageURL)
			_ = os.Remove(oldPath)
		}

		_ = os.MkdirAll("storages", os.ModePerm)

		ext := filepath.Ext(file.Filename)
		if ext == "" {
			ext = ".jpg"
		}

		slug := createSlug(targetTitle)
		filename := fmt.Sprintf("%s-%d%s", slug, time.Now().Unix(), ext)

		fields["image_url"] = filename
	}

	if req.CreatedBy != nil {
		fields["created_by"] = *req.CreatedBy
	}

	if req.ModifiedBy != nil && *req.ModifiedBy != "" {
		fields["modified_by"] = *req.ModifiedBy
	} else if authUser != "" {
		fields["modified_by"] = authUser
	}

	if len(fields) > 0 {
		if err := s.Repo.UpdatePartial(id, fields); err != nil {
			return nil, err
		}
	}

	return s.Repo.FindByID(id)
}

func (s *BookService) Delete(id int) error {
	existingBook, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	// Hapus file fisik dari storage
	if existingBook.ImageURL != nil && *existingBook.ImageURL != "" {
		filePath := filepath.Join("storages", *existingBook.ImageURL)
		_ = os.Remove(filePath)
	}

	return s.Repo.Delete(id)
}
