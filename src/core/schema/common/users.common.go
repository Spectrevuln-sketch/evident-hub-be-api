package common

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	Base
	Fullname              string     `gorm:"type:varchar(256);not null" json:"fullname"`
	Name                  string     `gorm:"type:varchar(256);not null" json:"name"`
	Username              string     `gorm:"type:varchar(256);not null;unique" json:"username"`
	Password              string     `gorm:"type:varchar(256);not null" json:"-"`
	RoleID                *uuid.UUID `gorm:"type:uuid;column:role_id; default:null" json:"role_id,omitempty"`
	Role                  *Roles     `gorm:"Foreignkey:RoleID;association_foreignkey:ID" json:"roles,omitempty"`
	IsTempPassword        bool       `gorm:"type:bool" json:"is_temp_password"`
	TempPasswordExpiresAt *time.Time `gorm:"type:date" json:"temp_password_expires_at"`
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if len(u.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	return nil
}
