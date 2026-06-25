package evidentPhoto

import (
	"github.com/google/uuid"
)

const (
	PHOTO_AFTER  = "after"
	PHOTO_BEFORE = "before"
)

type EvidentPhoto struct {
	Base

	EvidentID *uuid.UUID `gorm:"type:uuid;not null" json:"evident_id"`

	PhotoType string `gorm:"type:varchar(255);not null" json:"photo_type"`

	FilePath string `gorm:"type:varchar(255);not null" json:"file_path"`
}
