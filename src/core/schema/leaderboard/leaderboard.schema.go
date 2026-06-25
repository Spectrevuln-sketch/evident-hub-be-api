package leaderboard

import (
	"evidence-hub-be/src/core/schema/common"
	"evidence-hub-be/src/core/schema/evident"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	PointCreateEvidentH1 = 30
	PointCreateEvidentH2 = 20
	PointCreateEvidentH3 = 10
)

type LeaderBoard struct {
	Base
	UserID      uuid.UUID        `gorm:"type:varchar(256);" json:"user_id"`
	User        common.Users     `gorm:"foreignKey:UserID;references:ID" json:"user"`
	EvidentID   *uuid.UUID       `gorm:"type:varchar(256);" json:"evident_id,omitempty"`
	Evident     *evident.Evident `gorm:"foreignKey:EvidentID;references:ID" json:"evident,omitempty"`
	Point       int64            `gorm:"type:int;not null" json:"point"`
	Description string           `gorm:"type:varchar(256);not null" json:"description"`
}

func (r *LeaderBoard) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func GetPointByCategory(category string) int64 {
	switch category {
	case "H1":
		return PointCreateEvidentH1
	case "H2":
		return PointCreateEvidentH2
	case "H3":
		return PointCreateEvidentH3
	default:
		return 0
	}
}
