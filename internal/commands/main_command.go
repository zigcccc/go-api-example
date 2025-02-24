package commands

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"go_api_example/internal/auth"
	"go_api_example/internal/controllers"
	"go_api_example/internal/database"
	"go_api_example/internal/config"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Start a rest web server",
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func init() {
	config.InitConfig()

	RootCmd.AddCommand(restCmd)
}