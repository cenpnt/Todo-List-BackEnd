package models

type UserResponse struct {
	ID       uint   `json:"ID"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Tasks    []Task `json:"tasks"`
}
