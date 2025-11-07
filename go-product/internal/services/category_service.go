package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/models"
	"github.com/iqmalr-pedia/go-product/internal/repositories"
)

type CategoryService interface {
	CreateCategory(req *request.CreateCategoryRequest) (*response.CategoryResponse, error)
	GetCategoryByID(id uint) (*response.CategoryResponse, error)
	GetCategoryBySlug(slug string) (*response.CategoryResponse, error)
	UpdateCategory(id uint, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error)
	DeleteCategory(id uint) error
	ListCategories(req *request.ListCategoriesRequest) (*response.CategoryListResponse, error)
	GetCategoryTree(req *request.GetCategoryTreeRequest) (*response.CategoryTreeListResponse, error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(req *request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	slug := generateSlug(req.Name)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	level := 1
	if req.Level > 0 {
		level = req.Level
	} else if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
		level = parent.Level + 1
	}

	sortOrder := 0
	if req.SortOrder > 0 {
		sortOrder = req.SortOrder
	}

	category := &models.Category{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Slug:      slug,
		Desc:      req.Desc,
		Icon:      req.Icon,
		ImageURL:  req.ImageURL,
		Level:     level,
		SortOrder: sortOrder,
		IsActive:  isActive,
	}

	err := s.categoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *categoryService) GetCategoryByID(id uint) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *categoryService) GetCategoryBySlug(slug string) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *categoryService) UpdateCategory(id uint, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		category.ParentID = req.ParentID
		if req.ParentID != nil {
			parent, err := s.categoryRepo.GetByID(*req.ParentID)
			if err != nil {
				return nil, fmt.Errorf("parent category not found: %w", err)
			}
			category.Level = parent.Level + 1
		} else {
			category.Level = 1
		}
	}

	if req.Name != "" {
		category.Name = req.Name
		category.Slug = generateSlug(req.Name)
	}

	if req.Desc != "" {
		category.Desc = req.Desc
	}

	if req.Icon != "" {
		category.Icon = req.Icon
	}

	if req.ImageURL != "" {
		category.ImageURL = req.ImageURL
	}

	if req.Level > 0 {
		category.Level = req.Level
	}

	if req.SortOrder >= 0 {
		category.SortOrder = req.SortOrder
	}

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	hasChildren, err := s.categoryRepo.HasChildren(id)
	if err != nil {
		return err
	}

	if hasChildren {
		return errors.New("cannot delete category with child categories")
	}

	hasProducts, err := s.categoryRepo.HasProducts(id)
	if err != nil {
		return err
	}

	if hasProducts {
		return errors.New("cannot delete category with associated products")
	}

	return s.categoryRepo.Delete(id)
}

func (s *categoryService) ListCategories(req *request.ListCategoriesRequest) (*response.CategoryListResponse, error) {
	categories, err := s.categoryRepo.List(req.ParentID, req.Level, req.IsActive, req.Sort, req.Order)
	if err != nil {
		return nil, err
	}

	var categoryResponses []response.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, *s.toCategoryResponse(&category))
	}

	return &response.CategoryListResponse{
		Data: categoryResponses,
	}, nil
}

func (s *categoryService) GetCategoryTree(req *request.GetCategoryTreeRequest) (*response.CategoryTreeListResponse, error) {
	categories, err := s.categoryRepo.GetTree(req.IsActive)
	if err != nil {
		return nil, err
	}

	var categoryTreeResponses []response.CategoryTreeResponse
	for _, category := range categories {
		categoryTreeResponses = append(categoryTreeResponses, *s.toCategoryTreeResponse(&category))
	}

	return &response.CategoryTreeListResponse{
		Data: categoryTreeResponses,
	}, nil
}

func (s *categoryService) toCategoryResponse(category *models.Category) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID:        category.ID,
		UUID:      category.UUID.String(),
		ParentID:  category.ParentID,
		Name:      category.Name,
		Slug:      category.Slug,
		Desc:      category.Desc,
		Icon:      category.Icon,
		ImageURL:  category.ImageURL,
		Level:     category.Level,
		SortOrder: category.SortOrder,
		IsActive:  category.IsActive,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func (s *categoryService) toCategoryTreeResponse(category *models.Category) *response.CategoryTreeResponse {
	resp := &response.CategoryTreeResponse{
		ID:        category.ID,
		UUID:      category.UUID.String(),
		ParentID:  category.ParentID,
		Name:      category.Name,
		Slug:      category.Slug,
		Desc:      category.Desc,
		Icon:      category.Icon,
		ImageURL:  category.ImageURL,
		Level:     category.Level,
		SortOrder: category.SortOrder,
		IsActive:  category.IsActive,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if len(category.Children) > 0 {
		for _, child := range category.Children {
			resp.Children = append(resp.Children, *s.toCategoryTreeResponse(&child))
		}
	}

	return resp
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	return result.String()
}
