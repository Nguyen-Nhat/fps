package errorz

import (
	"fmt"
)

// Error show on Excel
var (
	ErrDefault = "xảy ra lỗi"
	ErrNoSites = func(siteCode string) string {
		return fmt.Sprintf("không tìm thấy thông tin site %s", siteCode)
	}
	ErrNoSkus = func(sku ...string) string {
		return fmt.Sprintf("không tìm thấy thông tin sku %s", sku)
	}
)

// Error show on Debug
const (
	ErrMissingParameter     = "missing parameter"
	ErrNotEqualNumberParams = "not equal number params"
	ErrFunctionNoSupport    = "not support function"
)

// Internal error
const (
	ErrInternal = "internal server error"
)
