package common

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Roles struct {
	Base
	Name      string         `gorm:"type:varchar(256);uniqueIndex" json:"name"`
	Privilage datatypes.JSON `gorm:"type:jsonb" json:"privilage"`
}

func (r *Roles) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func isValidType(value string) bool {
	// for _, v := range RoleTypes {
	// 	if v == value {
	// 		return true
	// 	}
	// }
	return false
}
