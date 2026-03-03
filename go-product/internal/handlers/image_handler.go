package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-product/internal/dto/request"
	"github.com/iqmalr-pedia/go-product/internal/dto/response"
	"github.com/iqmalr-pedia/go-product/internal/services"
)

type ImageHandler struct {
	imageService services.ImageService
}

func NewImageHandler(imageService services.ImageService) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

func (h *ImageHandler) UploadImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "File is required"})
		return
	}

	altText := c.DefaultQuery("alt_text", "")
	sortOrder := 0
	if so := c.Query("sort_order"); so != "" {
		if parsed, err := strconv.Atoi(so); err == nil {
			sortOrder = parsed
		}
	}
	isPrimary := false
	if ip := c.Query("is_primary"); ip == "true" || ip == "1" {
		isPrimary = true
	}

	image, err := h.imageService.UploadImage(uint(productID), file, altText, sortOrder, isPrimary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, image)
}

func (h *ImageHandler) GetImages(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	images, err := h.imageService.GetImagesByProductID(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, images)
}

func (h *ImageHandler) UpdateImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid image ID"})
		return
	}

	var req request.UpdateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: err.Error()})
		return
	}

	image, err := h.imageService.UpdateImage(uint(productID), uint(imageID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

func (h *ImageHandler) DeleteImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid image ID"})
		return
	}

	if err := h.imageService.DeleteImage(uint(productID), uint(imageID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Product image deleted successfully"})
}

func (h *ImageHandler) SetPrimaryImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid product ID"})
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{Message: "Invalid image ID"})
		return
	}

	if err := h.imageService.SetPrimaryImage(uint(productID), uint(imageID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Primary product image updated successfully"})
}
