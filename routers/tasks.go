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

func EditTask(c *gin.Context) {
	id := c.Param("id")

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
        return
	}

	var task models.Task
	if err := initializers.DB.Preload("SubTasks", recursiveSubTaskPreload).First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
	}

	if task.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not the owner of this task"})
        return
	}

	var data models.EditTaskRequest
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
	}

	if err := initializers.DB.Model(&task).Updates(models.Task{
		Title:       data.Title,
		Description: data.Description,
		IsCompleted: data.IsCompleted,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	var updatedTask models.Task

	if err := initializers.DB.Preload("SubTasks", recursiveSubTaskPreload).First(&updatedTask, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated", "task": updatedTask})
}

func ToggleTask(c *gin.Context) {
    id := c.Param("id")

    userID, exist := c.Get("userID")
    if !exist {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
        return
    }

    var task models.Task
    if err := initializers.DB.Preload("SubTasks", recursiveSubTaskPreload).First(&task, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    if task.UserID != userID {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not the owner of this task"})
        return
    }

    var data models.ToggleTaskRequest
    if err := c.BindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    task.IsCompleted = data.IsCompleted

    if err := initializers.DB.Model(&task).Update("is_completed", data.IsCompleted).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task updated", "task": task})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
        return
	}

	var task models.Task
	if err := initializers.DB.Preload("SubTasks", recursiveSubTaskPreload).First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
	}

	if task.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not the owner of this task"})
        return
	}

	if err := initializers.DB.Delete(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
