package model

import (
	"evidence-hub-be/src/core/schema/evident"

	"github.com/google/uuid"
)

type SCreateEvidentRequest struct {
	DealerID      string `json:"dealer_id" binding:"required"`
	Category      string `json:"category" binding:"required"`
	CatatanTemuan string `json:"catatan_temuan" binding:"required"`
}

func (s *SCreateEvidentRequest) MapToSchema(m *evident.Evident, createdBy string) {

	m.DealerID = uuid.MustParse(s.DealerID)

	m.Category = s.Category
	if s.CatatanTemuan != "" {
		m.CatatanTemuan = s.CatatanTemuan
	}
	m.Status = evident.StatusActive
	m.CreatedBy = createdBy

}
