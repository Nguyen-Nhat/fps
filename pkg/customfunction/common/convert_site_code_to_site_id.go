package customFunc

import (
	"errors"
	"net/http"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/constants"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/helpers"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

const (
	allSite          = "ALL"
	siteIdForAllSite = 0
)

// ConvertSiteCode2SiteId ...
func ConvertSiteCode2SiteId(siteCode string, sellerId string) FuncResult {
	siteCode = strings.TrimSpace(siteCode)
	// 1. Check if all site
	if strings.ToUpper(siteCode) == allSite {
		return FuncResult{Result: siteIdForAllSite}
	}

	// 2. Call api
	siteInfoResp, err := callApiGetSites(siteCode, sellerId)
	if err != nil {
		return FuncResult{ErrorMessage: errorz.ErrDefault}
	}

	// 3. Get site id from response
	for _, site := range siteInfoResp {
		if utils.EqualsIgnoreCase(site.SellerSiteCode, siteCode) {
			return FuncResult{Result: site.Id}
		}
	}

	return FuncResult{ErrorMessage: errorz.ErrNoSites(siteCode)}
}

func callApiGetSites(siteCode string, sellerId string) ([]SiteInfo, error) {
	// 1. Prepare call api
	httpClient := helpers.InitHttpClient()
	reqHeader := map[string]string{"Content-Type": "application/json"}
	reqParams := map[string]string{"sellerId": sellerId, "siteName": siteCode}

	// 2. Call api
	httpStatus, resBody, err := utils.SendHTTPRequest[any, GetSiteResponse](httpClient, http.MethodGet, constants.UrlApiGetSites, reqHeader, reqParams, nil)
	if err != nil {
		logger.Errorf("failed to call %v, got error=%v, resBody=%+v", constants.UrlApiGetSites, err, resBody)
		return nil, err
	}
	if httpStatus != http.StatusOK {
		logger.Errorf("failed to call %v, got httpStatus=%d, resBody=%+v", constants.UrlApiGetSites, httpStatus, resBody)
		return nil, errors.New(errorz.ErrInternal)
	}

	// 3. Return data
	return resBody.Data, nil
}
