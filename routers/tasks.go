package routers

import (
	"net/http"

	"github.com/cenpnt/Todo-List-BackEnd/initializers"
	"github.com/cenpnt/Todo-List-BackEnd/models"
	"github.com/gin-gonic/gin"
)


func GetAllTasks(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
		return
	}

	var user models.User
	if err := initializers.DB.Preload("Tasks", "parent_task_id IS NULL").Preload("Tasks.SubTasks", recursiveSubTaskPreload).First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} 

	c.JSON(http.StatusOK, user.Tasks)
}

func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
		return
	}

	task.UserID = userID.(uint)

	if task.ParentTaskID != nil {
		var parentTask models.Task
		// Dereference the pointer address to get the actual parentTaskID (parentTaskID is stored as a pointer (address))
		if err := initializers.DB.First(&parentTask, *task.ParentTaskID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent task does not exist"})
			return
		}
	}

	if err := initializers.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "task created successfully",
		"task": task,
	})
}