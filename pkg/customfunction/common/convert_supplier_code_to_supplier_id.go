package customFunc

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/patrickmn/go-cache"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/constants"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/helpers"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// ConvertSupplierCode2SupplierId ...
func ConvertSupplierCode2SupplierId(sellerId string, supplierCode string) FuncResult {
	supplierCode = strings.TrimSpace(supplierCode)
	if supplierCode == "" {
		return FuncResult{ErrorMessage: errorz.ErrMissingParameter}
	}

	// 1. Check in cache
	supplierId, found := cacheStore.Get(getKeySupplier(sellerId, supplierCode))
	if found {
		return FuncResult{Result: supplierId.(int)}
	}

	// 2. Call api
	supplierInfoResp, err := callApiGetSuppliers(sellerId)
	if err != nil {
		return FuncResult{ErrorMessage: errorz.ErrDefault}
	}

	// 3. Save to cache
	for _, supplier := range supplierInfoResp {
		cacheStore.Set(getKeySupplier(sellerId, supplier.Code), supplier.Id, cache.DefaultExpiration)
	}

	// 4. Return data
	supplierId, found = cacheStore.Get(getKeySupplier(sellerId, supplierCode))
	if found {
		return FuncResult{Result: supplierId.(int)}
	}
	return FuncResult{ErrorMessage: errorz.ErrNoSupplier(supplierCode)}
}

func callApiGetSuppliers(sellerId string) ([]SupplierInfo, error) {
	// 1. Prepare call api
	httpClient := helpers.InitHttpClient()
	reqHeader := map[string]string{"Content-Type": "application/json"}
	reqParams := map[string]string{"sellerId": sellerId}

	// 2. Call api
	httpStatus, resBody, err := utils.SendHTTPRequest[any, GetSupplierResponse](httpClient, http.MethodGet, constants.UrlGetSuppliers, reqHeader, reqParams, nil)
	if err != nil {
		logger.Errorf("failed to call %v, got error=%v, resBody=%+v", constants.UrlGetSuppliers, err, resBody)
		return nil, err
	}
	if httpStatus != http.StatusOK {
		logger.Errorf("failed to call %v, got httpStatus=%d, resBody=%+v", constants.UrlGetSuppliers, httpStatus, resBody)
		return nil, errors.New(errorz.ErrInternal)
	}

	// 3. Return data
	return resBody.Data.Suppliers, nil
}

func getKeySupplier(sellerId string, supplierCode string) string {
	return fmt.Sprintf("%s_supplier_%s", sellerId, supplierCode)
}
