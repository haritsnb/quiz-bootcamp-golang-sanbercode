package services

import (
	"app/models"
	"app/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(r *repositories.UserRepository) *UserService {
	return &UserService{Repo: r}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.Repo.FindAll()
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.Repo.FindByID(id)
}

func (s *UserService) Create(req models.UserCreateRequest, authUser string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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

	user := &models.User{
		Username:   req.Username,
		Password:   string(hashedPassword),
		CreatedBy:  &creator,
		ModifiedBy: &modifier,
	}

	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) Update(id int, req models.UserUpdateRequest, authUser string) error {
	if _, err := s.Repo.FindByID(id); err != nil {
		return err
	}

	fields := make(map[string]interface{})

	if req.Username != nil && *req.Username != "" {
		fields["username"] = *req.Username
	}

	if req.Password != nil && *req.Password != "" {
		if len(*req.Password) < 6 {
			return errors.New("password minimal 6 karakter")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		fields["password"] = string(hashedPassword)
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

func (s *UserService) Delete(id int) error {
	return s.Repo.Delete(id)
}
