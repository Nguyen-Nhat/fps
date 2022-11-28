package common

import "math"

func PageCalculator(total int, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

func PageSizeCalculator(total int, page int, pageSize int) int {
	var expectedPageSize int
	if (page-1)*pageSize > total {
		expectedPageSize = 0
	} else if temp := page * pageSize; temp > total {
		expectedPageSize = total - temp + pageSize
	} else {
		expectedPageSize = pageSize
	}
	return expectedPageSize
}
