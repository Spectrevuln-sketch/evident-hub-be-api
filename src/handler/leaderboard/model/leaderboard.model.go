package model

import (
	"github.com/google/uuid"
)

type LeaderBoardResponse struct {
	Rank int64 `json:"rank"`

	UserID uuid.UUID `json:"user_id"`

	Fullname string `json:"fullname"`

	TotalPoint int64 `json:"total_point"`

	TotalFind int64 `json:"total_find"`
}
