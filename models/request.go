package models

type ToggleTaskRequest struct {
	IsCompleted bool `json:"is_completed"`
}

type EditTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}
