package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go_api_example/internal/database"
	"go_api_example/internal/models"
	"go_api_example/internal/services"
)

// CreateProduct - POST /products
func CreateProduct(c *gin.Context) {
	var product models.Product

	// Bind JSON to product struct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user email from token
	userID, _ := c.Get("userID")

	// Get user by email (assuming user is already in DB)
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Associate the product with the user
	product.UserID = userID.(uuid.UUID)

	// Save to DB
	database.DB.Create(&product)
	c.JSON(http.StatusCreated, product)
}

// GetProducts - GET /products
func GetProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

// GetProduct - GET /products/:id
func GetProduct(c *gin.Context) {
	product, err := services.GetProductByIDAndUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return the product
	c.JSON(http.StatusOK, product)
}

// UpdateProduct - PUT /products/:id
func UpdateProduct(c *gin.Context) {
	product, err := services.GetProductByIDAndUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Bind request body to update product
	var input models.UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update product fields
	database.DB.Model(&product).Updates(input)

	// Update product
	database.DB.Save(&product)
	c.JSON(http.StatusOK, product)
}

// DeleteProduct - DELETE /products/:id
func DeleteProduct(c *gin.Context) {
	product, err := services.GetProductByIDAndUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Delete product
	database.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
