package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqmalr-pedia/go-auth/v2/internal/repositories"
)

type InternalHandler struct {
	userRepo repositories.UserRepositoryInterface
}

func NewInternalHandler(userRepo repositories.UserRepositoryInterface) *InternalHandler {
	return &InternalHandler{userRepo: userRepo}
}

type UserValidationResponse struct {
	ID    uint   `json:"id"`
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *InternalHandler) GetUserByEmailInternal(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email query parameter is required"})
		return
	}

	user, err := h.userRepo.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := UserValidationResponse{
		ID:    user.ID,
		UUID:  user.UUID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, response)
}

func (h *InternalHandler) GetUserByIDInternal(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.FindUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := UserValidationResponse{
		ID:    user.ID,
		UUID:  user.UUID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, response)
}
