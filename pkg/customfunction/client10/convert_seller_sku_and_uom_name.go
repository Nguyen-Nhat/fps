package funcClient10

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/constants"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/helpers"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
)

const (
	pageSizeKey    = "pageSize"
	pageSizeGetSku = "100"
	sellerSkusKey  = "sellerSkus"
)

// ConvertSellerSkus ...
func ConvertSellerSkus(jsonItems string, sellerId string) customFunc.FuncResult {
	// 1. Parse input
	var inputItems []ItemInput
	if err := json.Unmarshal([]byte(jsonItems), &inputItems); err != nil {
		return customFunc.FuncResult{ErrorMessage: errorz.ErrDefault}
	}
	if len(inputItems) == 0 {
		return customFunc.FuncResult{} // return (nil, "") value
	}

	// 2. Call api
	products, err := utils.BatchExecutingReturn(constants.BatchSizeQuerySku, inputItems, callApiGetSkus, sellerId)
	//products, err := callApiGetSkus(inputItems)
	logger.Infof("==== products=%+v, err=%v", products, err)
	if err != nil {
		return customFunc.FuncResult{ErrorMessage: errorz.ErrDefault}
	}
	if len(products) == 0 {
		return customFunc.FuncResult{ErrorMessage: errorz.ErrNoSkus()}
	}

	// 3. Convert response
	var outputItems []ItemOutput
	for _, inputItem := range inputItems {
		existed := false
		for _, product := range products {
			productSellerId := fmt.Sprintf("%d", product.SellerId)
			if utils.EqualsIgnoreCase(inputItem.SellerSku, product.SellerSku) &&
				utils.EqualsIgnoreCase(inputItem.UomName, product.UomName) &&
				productSellerId == sellerId { // because api still return products that belong to other sellerIds, so we have to check sellerId -> todo remove

				itemOutput := ItemOutput{product.Sku, inputItem.Quantity}
				outputItems = append(outputItems, itemOutput)
				existed = true
				break
			}
		}
		if !existed {
			return customFunc.FuncResult{ErrorMessage: errorz.ErrNoSkuWithUomName(inputItem.SellerSku, inputItem.UomName)}
		}
	}

	return customFunc.FuncResult{Result: outputItems}
}

func callApiGetSkus(subItems []ItemInput, sellerIds ...interface{}) ([]Product, error) {
	// 1. Convert input to param
	sellerSkus := converter.Map(subItems, func(i ItemInput) string { return i.SellerSku })
	sellerSkusStr := strings.Join(sellerSkus[:], ",")

	// 2. Prepare call api
	httpClient := helpers.InitHttpClient()
	reqHeader := map[string]string{"Content-Type": "application/json"}
	reqParams := map[string]string{sellerSkusKey: sellerSkusStr, pageSizeKey: pageSizeGetSku}
	if len(sellerIds) > 0 {
		sellerId := sellerIds[0]
		if reflect.TypeOf(sellerId).Kind() == reflect.String {
			reqParams["sellerIds"] = sellerId.(string)
		}
	}

	// 3. Call api
	httpStatus, resBody, err := utils.SendHTTPRequest[any, GetSkuResponse](httpClient, http.MethodGet, constants.UrlApiGetSkus, reqHeader, reqParams, nil)
	if err != nil || httpStatus != http.StatusOK {
		logger.Errorf("failed to call %v, got error=%v, httpStatus=%d, resBody=%+v", constants.UrlApiGetSkus, err, httpStatus, resBody)
		return nil, err
	}
	logger.Infof("==== call %v, httpStatus=%d, resBody=%+v", constants.UrlApiGetSkus, httpStatus, resBody)

	// 4. Return data
	return resBody.Data.Products, nil
}

// ConvertSellerSkuAndUomName2Sku ...
func ConvertSellerSkuAndUomName2Sku(sellerId, sellerSku, uomName string) customFunc.FuncResult {
	// 1. Call api
	products, err := callApiGetSkus([]ItemInput{{SellerSku: sellerSku, UomName: uomName}}, sellerId)
	if err != nil {
		return customFunc.FuncResult{ErrorMessage: errorz.ErrDefault}
	}
	if len(products) == 0 {
		return customFunc.FuncResult{ErrorMessage: errorz.ErrNoSkus(sellerSku)}
	}

	// 2. Convert response
	for _, product := range products {
		productSellerId := fmt.Sprintf("%d", product.SellerId)
		if utils.EqualsIgnoreCase(sellerSku, product.SellerSku) && utils.EqualsIgnoreCase(uomName, product.UomName) && productSellerId == sellerId {
			return customFunc.FuncResult{Result: product.Sku}
		}
	}

	return customFunc.FuncResult{ErrorMessage: errorz.ErrNoSkuWithUomName(sellerSku, uomName)}
}
