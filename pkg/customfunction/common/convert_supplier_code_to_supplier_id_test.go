package customFunc

import (
	"reflect"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

func Test_ConvertSupplierCodes2SupplierIds(t *testing.T) {
	// mock cache
	cacheStore = cache.New(15*time.Minute, 120*time.Minute)
	cacheStore.Set(getKeySupplier("1", "code_12"), 12, cache.DefaultExpiration)

	defer func() {
		cacheStore.Flush()
	}()

	type args struct {
		sellerId      string
		inputSupplier string
	}
	tests := []struct {
		name string
		args args
		want FuncResult
	}{
		{"test ConvertSupplierCode2SupplierId with supplier exist",
			args{"1", "code_12"},
			FuncResult{Result: 12, ErrorMessage: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertSupplierCode2SupplierId(tt.args.sellerId, tt.args.inputSupplier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertSupplierCode2SupplierId() = %v, want %v", got, tt.want)
			}
		})
	}
}
