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

	// Omit the password for security
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": user,
	})
}

func GetUsers(c *gin.Context) {
	var users models.User

	if err := initializers.DB.Preload("Tasks").Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}