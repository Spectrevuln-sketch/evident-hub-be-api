package commands

import (
	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/config/envs"
	"evidence-hub-be/src/core/schema/common"
	"log"

	"gorm.io/gorm"
)

func CreateUserAdmin(cfg *envs.Config) {
	var role common.Roles
	password := cfg.AppConfig.PassAdmin

	// Fetch the role
	if err := config.DB.First(&role, "name = ?", "AUDIT").Error; err != nil {
		log.Printf("Cannot Get Role: %s", err.Error())
		return
	}

	user := common.Users{
		Fullname: "audit",
		Username: "audit",
		Password: password,
		RoleID:   &role.ID,
	}

	err := config.DB.Where("username = ?", "audit").First(&user).Error
	if err != nil {
		// Only create the user if not found
		if err == gorm.ErrRecordNotFound {
			if createErr := config.DB.Create(&user).Error; createErr != nil {
				log.Printf("Error creating user: %s", createErr.Error())
			} else {
				log.Println("Successfully created user.")
			}
		} else {
			log.Printf("Error finding user: %s", err.Error())
		}
	} else {
		log.Println("User already exists.")
	}
}
