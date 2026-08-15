package services

import (
	"app/models"
	"app/repositories"
)

type CategoryService struct {
	Repo *repositories.CategoryRepository
}

func NewCategoryService(r *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{Repo: r}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.Repo.FindAll()
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.Repo.FindByID(id)
}

func (s *CategoryService) Create(req models.CategoryCreateRequest, authUser string) (*models.Category, error) {
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

	cat := &models.Category{
		Name:       req.Name,
		CreatedBy:  &creator,
		ModifiedBy: &modifier,
	}

	err := s.Repo.Create(cat)
	return cat, err
}

func (s *CategoryService) Update(id int, req models.CategoryUpdateRequest, authUser string) error {
	if _, err := s.Repo.FindByID(id); err != nil {
		return err
	}

	fields := make(map[string]interface{})

	if req.Name != nil && *req.Name != "" {
		fields["name"] = *req.Name
	}

	if req.CreatedBy != nil {
		fields["created_by"] = *req.CreatedBy
	}

	if req.ModifiedBy != nil && *req.ModifiedBy != "" {
		fields["modified_by"] = *req.ModifiedBy
	} else if authUser != "" {
		fields["modified_by"] = authUser
	}

	if len(fields) == 0 {
		return nil
	}

	return s.Repo.UpdatePartial(id, fields)
}

func (s *CategoryService) Delete(id int) error {
	return s.Repo.Delete(id)
}
