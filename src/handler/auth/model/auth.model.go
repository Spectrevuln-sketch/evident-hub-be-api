package model

import (
	"evidence-hub-be/src/core/schema/common"
	"time"

	"github.com/google/uuid"
)

type SPayloadRegister struct {
	Fullname string `json:"fullname" example:"fullname"`
	Password string `json:"password" example:"password"`
	Username string `json:"username" example:"username"`
	RoleID   string `json:"role_id" example:"uuid-role"`
}

type SPayloadLogin struct {
	Username string `json:"username" example:"username"`
	Password string `json:"password" example:"password"`
}

type SPayloadChangePassword struct {
	Username    string `json:"username" example:"username"`
	OldPassword string `json:"old_password" example:"old_password"`
	NewPassword string `json:"new_password" example:"new_password"`
}

type UserResponse struct {
	common.Users
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (s *SPayloadRegister) MapToSchema(m *common.Users, roleId *uuid.UUID) {

	if s.Fullname != "" {
		m.Username = s.Username
		m.Fullname = s.Fullname
	}

	if s.Password != "" {
		m.Password = s.Password
	}

	m.RoleID = roleId

	if m.Password == "" {
		expiresAt := time.Now().Add((time.Hour * 24) * 10)

		m.TempPasswordExpiresAt = &expiresAt
		m.IsTempPassword = true
	}
}
