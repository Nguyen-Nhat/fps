package customFunc

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
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

var cacheStore = cache.New(15*time.Minute, 120*time.Minute)

// ConvertSiteCode2SiteId ...
func ConvertSiteCode2SiteId(siteCode string, sellerId string, useCache ...bool) FuncResult {
	siteCode = strings.TrimSpace(siteCode)
	// 1. Check if all site
	if strings.ToUpper(siteCode) == allSite {
		return FuncResult{Result: siteIdForAllSite}
	}

	// 2. Check in cache
	isUseCache := len(useCache) > 0 && useCache[0]
	if isUseCache {
		siteId, found := cacheStore.Get(getKeySite(sellerId, siteCode))
		if found {
			return FuncResult{Result: siteId.(int)}
		}
	}

	// 3. Call api
	siteInfoResp, err := callApiGetSites(siteCode, sellerId)
	if err != nil {
		return FuncResult{ErrorMessage: errorz.ErrDefault}
	}

	// 4. Get site id from response
	for _, site := range siteInfoResp {
		// 4.1. Save to cache
		if isUseCache {
			cacheStore.Set(getKeySite(sellerId, site.SellerSiteCode), site.Id, cache.DefaultExpiration)
		}
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

// ConvertSiteCodes2SiteIds ...
func ConvertSiteCodes2SiteIds(sellerId string, inputSiteCodes string, separator string) FuncResult {
	if inputSiteCodes == constant.EmptyString || inputSiteCodes == allSite {
		return ConvertSiteCode2SiteId(allSite, sellerId)
	}

	if separator == constant.EmptyString {
		separator = constant.SplitByNewLine
	}
	siteCodes := strings.Split(inputSiteCodes, separator)
	var listSiteIds []int
	for _, siteCode := range siteCodes {
		if siteCode == constant.EmptyString {
			continue
		}
		siteId := ConvertSiteCode2SiteId(siteCode, sellerId, true)
		if siteId.ErrorMessage != constant.EmptyString {
			return FuncResult{ErrorMessage: siteId.ErrorMessage}
		}
		listSiteIds = append(listSiteIds, siteId.Result.(int))
	}
	return FuncResult{Result: listSiteIds}
}

func getKeySite(sellerId string, siteCode string) string {
	return fmt.Sprintf("%s_site_%s", sellerId, siteCode)
}
