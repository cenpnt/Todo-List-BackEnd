package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Tasks    []Task `json:"tasks"`
}

type Task struct {
	gorm.Model
	Title        string `gorm:"not null"`
	Description  string
	UserID       uint
	ParentTaskID *uint    `gorm:"index"`
	ParentTask   *Task    `gorm:"foreignKey:ParentTaskID"`
	SubTasks     []Task   `gorm:"foreignKey:ParentTaskID"`
	IsCompleted  bool     `gorm:"default:false"`
}