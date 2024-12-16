package routers

import (
	"net/http"

	"github.com/cenpnt/Todo-List-BackEnd/initializers"
	"github.com/cenpnt/Todo-List-BackEnd/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request data"})
		return
	}

	var existingUser models.User
	if err := initializers.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Email already exist"})
		return
	}

	if err := initializers.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Username already exist"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Could not hash password"})
		return
	}

	user.Password = string(hashedPassword)

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to register user"})
		return
	}

	userResponse := models.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": userResponse,
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := initializers.DB.Preload("Tasks").Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to get users"})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID: user.ID,
			Email: user.Email,
			Username: user.Username,
			Tasks: user.Tasks,
		})
	}

	c.JSON(http.StatusOK, userResponses)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := initializers.DB.Preload("Tasks").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse := models.UserResponse {
		ID: user.ID,
		Email: user.Email,
		Username: user.Username,
	}

	c.JSON(http.StatusOK, userResponse)
}