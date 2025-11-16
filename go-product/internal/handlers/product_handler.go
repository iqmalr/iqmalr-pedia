package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body request.CreateProductRequest true "Product data"
// @Success 201 {object} response.ProductResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ProductDetailResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductBySlug godoc
// @Summary Get product by slug
// @Description Get product by slug
// @Tags products
// @Produce json
// @Param slug path string true "Product slug"
// @Success 200 {object} response.ProductDetailResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products/slug/{slug} [get]
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Slug is required"})
		return
	}

	product, err := h.productService.GetProductBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body request.UpdateProductRequest true "Product data"
// @Success 200 {object} response.ProductResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	product, err := h.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	err = h.productService.DeleteProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Product deleted successfully"})
}

// ListProducts godoc
// @Summary List products
// @Description List products with optional filters
// @Tags products
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search term"
// @Param vendor_id query int false "Vendor ID"
// @Param category_id query int false "Category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param status query string false "Product status" Enums(draft, pending, active, rejected, out_of_stock)
// @Param is_featured query bool false "Is featured"
// @Param is_published query bool false "Is published"
// @Param sort query string false "Sort by" Enums(name, price, created_at, updated_at, view_count, sold_count, rating_avg)
// @Param order query string false "Sort order" Enums(asc, desc)
// @Success 200 {object} response.ProductListResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var req request.ListProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	products, err := h.productService.ListProducts(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProductStatus godoc
// @Summary Update product status (admin only)
// @Description Update product status (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param status body request.UpdateProductStatusRequest true "Status"
// @Success 200 {object} response.ProductStatusResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products/{id}/status [put]
func (h *ProductHandler) UpdateProductStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	var req request.UpdateProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	status, err := h.productService.UpdateProductStatus(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// PublishProduct godoc
// @Summary Publish product
// @Description Publish product
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ProductPublishResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products/{id}/publish [put]
func (h *ProductHandler) PublishProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	result, err := h.productService.PublishProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UnpublishProduct godoc
// @Summary Unpublish product
// @Description Unpublish product
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ProductPublishResponse
// @Failure 400 {object} response.MessageResponse
// @Failure 401 {object} response.MessageResponse
// @Failure 403 {object} response.MessageResponse
// @Failure 404 {object} response.MessageResponse
// @Failure 500 {object} response.MessageResponse
// @Router /products/{id}/unpublish [put]
func (h *ProductHandler) UnpublishProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	result, err := h.productService.UnpublishProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
