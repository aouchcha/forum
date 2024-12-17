package helpers

import "math"

func Pagination(page, DBlength int) int {
	var offset int
	if page < 1 || page > int(math.Ceil(float64(DBlength/10)+1)) {
		if page < 1 {
			page = 1
		} else {
			page = int(math.Ceil(float64(DBlength/10) + 1))
		}
		offset = DBlength - (DBlength - (10 * (page - 1)))
	} else {
		offset = DBlength - (DBlength - (10 * (page - 1)))
	}
	return offset
}
