package customFunc

import (
	"reflect"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

func Test_GetSites(t *testing.T) {
	// happy case: all site
	res := ConvertSiteCode2SiteId("all", "12")
	logger.Infof("Result = %+v", res)

	// happy case: specific site
	res3 := ConvertSiteCode2SiteId("S1000", "12")
	logger.Infof("Result = %+v", res3)

	// error case: not exist
	res2 := ConvertSiteCode2SiteId("SBN1000", "12")
	logger.Infof("Result = %+v", res2)
}

func Test_ConvertSiteCodes2SiteIds(t *testing.T) {
	// mock cache
	cacheStore = cache.New(15*time.Minute, 120*time.Minute)
	cacheStore.Set(getKeySite("1", "code_12"), 12, cache.DefaultExpiration)
	cacheStore.Set(getKeySite("1", "code_13"), 13, cache.DefaultExpiration)
	cacheStore.Set(getKeySite("1", "code_14"), 14, cache.DefaultExpiration)

	defer func() {
		cacheStore.Flush()
	}()

	type args struct {
		sellerId       string
		inputSiteCodes string
		separator      string
	}
	tests := []struct {
		name string
		args args
		want FuncResult
	}{
		{"test ConvertSiteCodes2SiteIds with list site default separate",
			args{"1", "code_12\ncode_13", ""},
			FuncResult{Result: []int{12, 13}, ErrorMessage: ""},
		},
		{"test ConvertSiteCodes2SiteIds with site with separate",
			args{"1", "code_13\ncode_14", "\n"},
			FuncResult{Result: []int{13, 14}, ErrorMessage: ""},
		},
		{"test ConvertSiteCodes2SiteIds with site with duplicate separate",
			args{"1", "code_13\n\n\ncode_14", "\n"},
			FuncResult{Result: []int{13, 14}, ErrorMessage: ""},
		},
		{"test ConvertSiteCodes2SiteIds with site with duplicate end separate",
			args{"1", "\ncode_13\ncode_14\n\n\n", "\n"},
			FuncResult{Result: []int{13, 14}, ErrorMessage: ""},
		},
		{"test ConvertSiteCodes2SiteIds with all site",
			args{"1", "", "all"},
			FuncResult{Result: 0, ErrorMessage: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertSiteCodes2SiteIds(tt.args.sellerId, tt.args.inputSiteCodes, tt.args.separator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertSiteCodes2SiteIds() = %v, want %v", got, tt.want)
			}
		})
	}
}
