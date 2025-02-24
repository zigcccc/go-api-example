package services

import (
	"errors"
	"go_api_example/internal/database"
	"go_api_example/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProductByIDAndUser(c *gin.Context) (*models.Product, error) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("Unauthorized")
	}

	// Get product ID from request params
	productID := c.Param("id")

	// Fetch product from database (ONLY IF USER IS OWNER)
	var product models.Product
	result := database.DB.Where("id = ? AND user_id = ?", productID, userID).First(&product)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, errors.New("NotFound")
	} else if result.Error != nil {
		return nil, result.Error
	}

	// RETURN STRONK PRODUCT 💪
	return &product, nil
}
