package customFunc

import (
	"testing"

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
