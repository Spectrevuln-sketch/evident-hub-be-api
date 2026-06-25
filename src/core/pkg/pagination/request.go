package pagination

import (
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	Page   int
	Limit  int
	Offset int
}

func BindQuery(c *fiber.Ctx) (*Request, error) {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	return &Request{
		Limit:  limit,
		Page:   page,
		Offset: offset,
	}, nil
}
