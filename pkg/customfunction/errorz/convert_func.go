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

	ErrCallAPI = func(api, err string) string {
		return fmt.Sprintf("Call API %s thất bại. Lỗi: %s", api, err)
	}
	ErrNoSkuWithUomName = func(sellerSku, uomName string) string {
		return fmt.Sprintf("Không tìm thấy thông tin sku với sellerSku=%s, uomName=%s", sellerSku, uomName)
	}
	ErrDateTimeFormat = func(dateTime, format string) string {
		return fmt.Sprintf("Sai định dạng ngày tháng. Ngày tháng: %s, định dạng: %s", dateTime, format)
	}
	ErrCantParseValue = func(value, valueType string) string {
		return fmt.Sprintf("Can not parse value %s to %s", value, valueType)
	}
)

// Error show on Debug
const (
	ErrMissingParameter     = "missing parameter"
	ErrNotEqualNumberParams = "not equal number params"
	ErrLessThanNumberParams = "less than number params"
	ErrFunctionNoSupport    = "not support function"
)

// Internal error
const (
	ErrInternal            = "internal server error"
	ErrSellerIdIsNotNumber = "sellerId is not a number"
)
