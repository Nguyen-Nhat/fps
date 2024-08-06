package customFunc

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/constants"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/helpers"
	t "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customtype"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

const (
	allSite          = "ALL"
	siteIdForAllSite = 0

	reqParamSellerId = "sellerId"
	reqParamSiteName = "siteName"
	reqParamSiteIds  = "siteIds"
	reqParamIsActive = "isActive"

	filterActiveSite = "true"

	apiName = "GetSites"
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

func callApiGetSites(siteCode string, sellerId string, siteIds ...string) ([]SiteInfo, error) {
	// 1. Prepare call api
	httpClient := helpers.InitHttpClient()
	reqHeader := map[string]string{"Content-Type": "application/json"}
	reqParams := []t.Pair[string, string]{
		{Key: reqParamSellerId, Value: sellerId},
		{Key: reqParamSiteName, Value: siteCode}, // Filter results like Site Name or Site Code.
	}

	if len(siteIds) > 0 {
		for _, siteId := range siteIds {
			reqParams = append(reqParams, t.Pair[string, string]{Key: reqParamSiteIds, Value: siteId})
		}
	}

	// 1.1. Check if client enable for OMNI-1139
	sellerIdInt, err := strconv.Atoi(sellerId)
	if err != nil {
		logger.Errorf("failed to convert sellerId=%s to int, got error=%v", sellerId, err)
		return nil, errors.New(errorz.ErrSellerIdIsNotNumber)
	}
	if utils.Contains(config.Cfg.ExtraConfig.Epic1139EnableSellersObj, int32(sellerIdInt)) {
		reqParams = append(reqParams, t.Pair[string, string]{Key: reqParamIsActive, Value: filterActiveSite}) // OMNI-1139, get only active site
	}

	// 2. Call api
	httpStatus, resBody, err := utils.SendHTTPRequestWithArrayParams[any, GetSiteResponse](httpClient, http.MethodGet, constants.UrlApiGetSites, reqHeader, reqParams, nil)
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
	if inputSiteCodes == constant.EmptyString || strings.ToUpper(inputSiteCodes) == allSite {
		return FuncResult{}
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

// ValidateAndConvertSiteCode2SiteId
// Validate siteCode have siteId in list mustInSiteIds, if true return siteId, else return error
// Req: sellerId, siteCode, mustInSiteIds (can be array string siteId or only a string siteId)
// Res: siteId if siteCode have siteId in mustInSiteIds, else return error
func ValidateAndConvertSiteCode2SiteId(sellerId, siteCode string, mustInSiteIds interface{}) FuncResult {
	siteCode = strings.TrimSpace(siteCode)
	// 1. Check if all site
	if strings.ToUpper(siteCode) == allSite {
		return FuncResult{Result: siteIdForAllSite}
	}

	if mustInSiteIds == nil {
		return ConvertSiteCode2SiteId(siteCode, sellerId, true)
	}

	// parse siteIds
	siteIdsArr := make([]string, 0)
	siteIdParsed := gjson.Parse(mustInSiteIds.(string))
	if siteIdParsed.IsArray() {
		for _, siteId := range siteIdParsed.Array() {
			siteIdsArr = append(siteIdsArr, siteId.String())
		}
	} else {
		siteIdsArr = append(siteIdsArr, mustInSiteIds.(string))
	}

	if len(siteIdsArr) == 0 {
		return ConvertSiteCode2SiteId(siteCode, sellerId, true)
	}

	// 2. Check in cache
	siteId, found := cacheStore.Get(getKeySite(sellerId, siteCode))
	if found {
		for _, id := range siteIdsArr {
			if id == strconv.Itoa(siteId.(int)) {
				return FuncResult{Result: siteId.(int)}
			}
		}
		return FuncResult{ErrorMessage: errorz.ErrNoSites(siteCode)}
	}

	// 3. Call api
	siteInfoResp, err := callApiGetSites(siteCode, sellerId, siteIdsArr...)
	if err != nil {
		return FuncResult{ErrorMessage: errorz.ErrCallAPI(apiName, err.Error())}
	}

	// 4. Get site id from response
	for _, site := range siteInfoResp {
		// 4.1. Save to cache
		cacheStore.Set(getKeySite(sellerId, site.SellerSiteCode), site.Id, cache.DefaultExpiration)
		if utils.EqualsIgnoreCase(site.SellerSiteCode, siteCode) {
			return FuncResult{Result: site.Id}
		}
	}

	return FuncResult{ErrorMessage: errorz.ErrNoSites(siteCode)}
}
