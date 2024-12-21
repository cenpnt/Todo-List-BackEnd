package routers

import (
	"net/http"

	"github.com/cenpnt/Todo-List-BackEnd/initializers"
	"github.com/cenpnt/Todo-List-BackEnd/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskWithProgress struct {
    models.Task
    Progress float64 `json:"progress"`
}

func calculateTaskProgress(task *models.Task, db *gorm.DB) float64 {
    if len(task.SubTasks) == 0 {
        if task.IsCompleted {
            return 100.0
        }
        return 0.0
    }

    var totalProgress float64 = 0.0
    allSubtasksCompleted := true

    for _, subtask := range task.SubTasks {
        progress := calculateTaskProgress(&subtask, db)
        totalProgress += progress
        
        if !subtask.IsCompleted {
            allSubtasksCompleted = false
        }
    }

    if allSubtasksCompleted && !task.IsCompleted {
		// If all of the subtasks is completed but the parent task is not completed
        task.IsCompleted = true
        db.Model(&models.Task{}).Where("id = ?", task.ID).Update("is_completed", true)
    } else if !allSubtasksCompleted && task.IsCompleted {
		// If all of the subtasks is not completed but the parent task is completed
		task.IsCompleted = false
		db.Model(&models.Task{}).Where("id = ?", task.ID).Update("is_completed", false)
	}
    
    return totalProgress / float64(len(task.SubTasks))
}

func convertToTaskWithProgress(tasks []models.Task, db *gorm.DB) []TaskWithProgress {
    var result []TaskWithProgress
    for _, task := range tasks {
        taskWithProgress := TaskWithProgress{
            Task:     task,
            Progress: calculateTaskProgress(&task, db),
        }
        result = append(result, taskWithProgress)
    }
    return result
}

func GetAllTasks(c *gin.Context) {
    userID, exist := c.Get("userID")
    if !exist {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
        return
    }

    var user models.User
    if err := initializers.DB.Preload("Tasks", "parent_task_id IS NULL").
        Preload("Tasks.SubTasks", recursiveSubTaskPreload).
        First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    tasksWithProgress := convertToTaskWithProgress(user.Tasks, initializers.DB)
    c.JSON(http.StatusOK, tasksWithProgress)
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

	if task.ParentTaskID != nil {
		var parentTask models.Task
		if err := initializers.DB.Preload("SubTasks", recursiveSubTaskPreload).First(&parentTask, task.ParentTaskID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch parent task"})
            return
		}
		calculateTaskProgress(&parentTask, initializers.DB)
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
