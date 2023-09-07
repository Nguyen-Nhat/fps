package request

import (
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
)

func TestPageRequest_InitDefaultValue(t *testing.T) {
	type fields struct {
		Page     int
		PageSize int
	}
	tests := []struct {
		name            string
		fields          fields
		wantPageRequest PageRequest
	}{
		{"test normal case -> keep original value", fields{10, 120},
			PageRequest{10, 120}},
		{"test case equal default -> keep original value", fields{1, 10},
			PageRequest{1, 10}},
		{"test case equal default -> keep original value", fields{1, 10},
			PageRequest{1, 10}},
		{"test case page is negative -> keep original value", fields{-1, 20},
			PageRequest{-1, 20}},
		{"test case page is zero -> set default pageSize", fields{0, 20},
			PageRequest{constant.PaginationDefaultPage, 20}},
		{"test case pageSize is negative -> keep original value", fields{2, -10},
			PageRequest{2, -10}},
		{"test case pageSize is zero -> set default pageSize", fields{2, 0},
			PageRequest{2, constant.PaginationDefaultSize}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := PageRequest{
				Page:     tt.fields.Page,
				PageSize: tt.fields.PageSize,
			}

			c.InitDefaultValue()

			if !reflect.DeepEqual(c, tt.wantPageRequest) {
				t.Errorf("InitDefaultValue() got = %v, want %v", c, tt.wantPageRequest)
			}
		})
	}
}

func TestPageRequest_ValidatePagination(t *testing.T) {
	type fields struct {
		Page     int
		PageSize int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"test normal case", fields{2, 120}, false},

		{"test case page is negative", fields{-2, 120}, true},
		{"test case page is zero", fields{0, 10}, false},
		{"test case page is PaginationMaxPage", fields{constant.PaginationMaxPage, 10}, false},
		{"test case page is greater than PaginationMaxPage", fields{constant.PaginationMaxPage + 1, 10}, true},

		{"test case pageSize is negative", fields{2, -10}, true},
		{"test case pageSize is zero", fields{2, 0}, false},
		{"test case pageSize is PaginationMaxSize", fields{2, constant.PaginationMaxSize}, false},
		{"test case pageSize is greater than PaginationMaxPage", fields{3, constant.PaginationMaxSize + 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PageRequest{
				Page:     tt.fields.Page,
				PageSize: tt.fields.PageSize,
			}
			if err := c.ValidatePagination(); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePagination() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
