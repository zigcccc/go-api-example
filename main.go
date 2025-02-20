package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"go_api_example/auth"
	"go_api_example/controllers"
	"go_api_example/database"
)

func main() {
	database.ConnectDatabase()

	router := gin.Default()

	authRoutes := router.Group("/auth")
	{
		authRoutes.GET("/login", auth.LoginHandler)
		authRoutes.GET("/callback", auth.CallbackHandler)
	}

	productRoutes := router.Group("/products")
	productRoutes.Use(auth.AuthMiddleware())
	{
		productRoutes.POST("", controllers.CreateProduct)
		productRoutes.GET("", controllers.GetProducts)
		productRoutes.GET("/:id", controllers.GetProduct)
		productRoutes.PATCH("/:id", controllers.UpdateProduct)
		productRoutes.DELETE("/:id", controllers.DeleteProduct)
	}

	// Example route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	log.Println("Starting server on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
