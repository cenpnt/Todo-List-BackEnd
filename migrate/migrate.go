package main

import (
	"log"

	"github.com/cenpnt/Todo-List-BackEnd/initializers"
	"github.com/cenpnt/Todo-List-BackEnd/models"
)

func main() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	// Run migrations
	if err := initializers.DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed successfully!")
}
