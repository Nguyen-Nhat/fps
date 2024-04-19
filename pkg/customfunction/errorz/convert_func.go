package errorz

import (
	"errors"
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
	ErrNoSupplier = func(supplierCode string) string {
		return fmt.Sprintf("không tìm thấy thông tin supplier %s", supplierCode)
	}

	ErrFormatOrderScheduler = errors.New("Cấu trúc lịch đặt không đúng! Cấu trúc đúng có dạng T2->6,CN hoặc N2->5,9->20,30. Các số theo thứ tự tăng dần. Ex: T2->6,CN; TCN; T8; T5->Cn; T2,4,5; N1->5,9->20,31")
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
