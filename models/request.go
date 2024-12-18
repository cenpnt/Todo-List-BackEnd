package models

type ToggleTaskRequest struct {
	IsCompleted bool `json:"is_completed"`
}
