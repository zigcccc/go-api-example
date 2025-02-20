package auth

import (
	"fmt"
	"time"

	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"go_api_example/database"
	"go_api_example/models"
)

func LoginHandler(c *gin.Context) {
	authUrl := OAuth2Config.AuthCodeURL(OAuth2StateString, oauth2.AccessTypeOffline)
	c.JSON(http.StatusFound, gin.H{"authUrl": authUrl})
}

func CallbackHandler(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Exchange the code for a token
	token, err := OAuth2Config.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get user info
	client := OAuth2Config.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var userData struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	var user models.User
	result := database.DB.Where("email = ?", userData.Email).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// User doesn't exist, create new one
		user = models.User{Email: userData.Email}
		database.DB.Create(&user)
		fmt.Println("New user created:", user.Email)
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	} else {
		fmt.Println("User already exists:", user.Email)
	}

	// Create JWT Token
	tokenString, err := CreateJWTToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"email":   user.Email,
		"token":   tokenString,
	})
}

func CreateJWTToken(id uuid.UUID, email string) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}
