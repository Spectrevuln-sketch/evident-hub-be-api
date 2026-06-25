package pagination

type Response struct {
	Page         int `json:"page"`
	TotalPerPage int `json:"total_per_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	Items        any `json:"items"`
}

func NewResponse(req *Request, data any, total int) *Response {
	return &Response{
		Page:         req.Page,
		TotalPerPage: req.Limit,
		TotalItems:   total,
		Items:        data,
		TotalPages:   calculateTotalPages(total, req.Limit),
	}
}
