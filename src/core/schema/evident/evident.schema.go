package evident

import (
	"evidence-hub-be/src/core/schema/common"
	"evidence-hub-be/src/core/schema/evidentPhoto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	CategoryH1 = "H1"
	CategoryH2 = "H2"
	CategoryH3 = "H3"
)

const (
	StatusActive    = "active"
	StatusCompleted = "completed"
)

type Evident struct {
	Base
	DealerID      uuid.UUID                   `gorm:"type:varchar(256);" json:"dealer_id"`
	Dealer        common.Users                `gorm:"foreignKey:DealerID;references:ID" json:"dealer"`
	CreatedBy     string                      `gorm:"type:varchar(256);not null" json:"created_by"`
	Category      string                      `gorm:"type:varchar(256);not null" json:"category"`
	CatatanTemuan string                      `gorm:"type:varchar(256)" json:"catatan_temuan"`
	Status        string                      `gorm:"type:varchar(256);not null" json:"status"`
	Photos        []evidentPhoto.EvidentPhoto `gorm:"foreignKey:EvidentID" json:"photos"`
}

func (r *Evident) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func isValidCategory(category string) bool {
	switch category {
	case CategoryH1,
		CategoryH2,
		CategoryH3:
		return true
	default:
		return false
	}
}
func isValidStatus(status string) bool {
	switch status {
	case StatusActive,
		StatusCompleted:
		return true
	default:
		return false
	}
}
