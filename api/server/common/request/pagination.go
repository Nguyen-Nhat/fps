package request

import (
	"fmt"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
)

// PageRequest ...
type PageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// InitDefaultValue ... init default value: page=1 and pageSize=10
func (c *PageRequest) InitDefaultValue() {
	if c.Page == 0 {
		c.Page = constant.PaginationDefaultPage
	}

	if c.PageSize == 0 {
		c.PageSize = constant.PaginationDefaultSize
	}
}

// ValidatePagination ... validate page and pageSize
func (c *PageRequest) ValidatePagination() error {
	if c.Page < 0 || c.Page > constant.PaginationMaxPage {
		return fmt.Errorf("page is out of range")
	}

	if c.PageSize < 0 || c.PageSize > constant.PaginationMaxSize {
		return fmt.Errorf("pageSize is out of range")
	}

	return nil
}

func (c *PageRequest) Offset() int {
	return (c.Page - 1) * c.PageSize
}

func (c *PageRequest) Limit() int {
	return c.PageSize
}
