package pagination

import "math"

func calculateTotalPages(totalItems, itemsPerPage int) int {
	if itemsPerPage <= 0 {
		return 0
	}
	return int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))
}
