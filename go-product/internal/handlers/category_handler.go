package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/services"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param category body request.CreateCategoryRequest true "Category data"
// @Success 201 {object} response.CategoryResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.CategoryResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid category ID"})
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryBySlug godoc
// @Summary Get category by slug
// @Description Get category by slug
// @Tags categories
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} response.CategoryResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Slug is required"})
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body request.UpdateCategoryRequest true "Category data"
// @Success 200 {object} response.CategoryResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid category ID"})
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete category (admin only)
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid category ID"})
		return
	}

	err = h.categoryService.DeleteCategory(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Category deleted successfully"})
}

// ListCategories godoc
// @Summary List categories
// @Description List categories with optional filters
// @Tags categories
// @Produce json
// @Param parent_id query int false "Parent ID"
// @Param level query int false "Category level"
// @Param is_active query bool false "Is active"
// @Param sort query string false "Sort by" Enums(sort_order, name, level, created_at, updated_at)
// @Param order query string false "Sort order" Enums(asc, desc)
// @Success 200 {object} response.CategoryListResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	var req request.ListCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	categories, err := h.categoryService.ListCategories(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryTree godoc
// @Summary Get category tree
// @Description Get category tree structure
// @Tags categories
// @Produce json
// @Param is_active query bool false "Is active"
// @Success 200 {object} response.CategoryTreeListResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /categories/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	var req request.GetCategoryTreeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	categories, err := h.categoryService.GetCategoryTree(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
