package commands

import (
	"evidence-hub-be/src/core/config"
	"evidence-hub-be/src/core/schema/common"
	"log"
	// "strings"
)

func CreateRoleAdmin() {
	// jsonMenu := `["DASHBOARD", "MERCHANTS", "TRANSACTIONS", "WITHDRAWAL", "SETTING", "PROFILE", "WHITELIST"]`

	var user common.Users
	user_audit := config.DB.Where("username = ?", "audit").First(&user)
	if user_audit.Error != nil {
		log.Printf("failed to retrieve user %s " + user_audit.Error.Error())
	}

	role := common.Roles{
		Name: "AUDIT",
		// Menu: jsonMenu,
		// UserID: user.ID,
	}
	err := config.DB.Where("name = ?", "AUDIT").First(&role).Error
	if err != nil {
		config.DB.Create(&role)
		log.Println("Succesfully create role.")
		return
	}

	user_dealer := config.DB.Where("username = ?", "dealer").First(&user)
	if user_dealer.Error != nil {
		log.Printf("failed to retrieve user %s " + user_dealer.Error.Error())
	}

	role_dealer := common.Roles{
		Name: "DEALER",
		// Menu: jsonMenu,
		// UserID: user.ID,
	}

	if err := config.DB.Where("name = ?", "DEALER").First(&role_dealer).Error; err != nil {
		config.DB.Create(&role_dealer)
		log.Println("Succesfully create role.")
		return
	}
}
