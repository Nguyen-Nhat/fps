package funcClient10

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/converter"
)

const FuncConvertSellerSkuAndUomName = "convertSellerSkuAndUomName"

var (
	errDefault = customFunc.FuncResult{ErrorMessage: "xảy ra lỗi"}
	errNoSkus  = customFunc.FuncResult{ErrorMessage: "không tìm thấy thông tin sku"}
)

const (
	// urlApiGetSkus ... has path is `/api/v2/skus`
	// api-doc: https://apidoc.teko.vn/project-doc/approved/core_logic_layer/retail/catalog/version/latest/paths/api-v2-skus/get
	urlApiGetSkus = "http://catalog-core-v2-api.catalog/api/v2/skus" // url call service name
	//urlApiGetSkus = "http://localhost:10080/api/v2/skus" // url call local

	batchSizeQuerySku = 50
)

// ConvertSellerSkus ...
func ConvertSellerSkus(jsonItems string) customFunc.FuncResult {
	// 1. Parse input
	var inputItems []ItemInput
	if err := json.Unmarshal([]byte(jsonItems), &inputItems); err != nil {
		return errDefault
	}
	if len(inputItems) == 0 {
		return customFunc.FuncResult{} // return (nil, "") value
	}

	// 2. Call api
	products, err := utils.BatchExecutingReturn(batchSizeQuerySku, inputItems, callApiGetSkus)
	//products, err := callApiGetSkus(inputItems)
	if err != nil {
		return errDefault
	} else if len(products) == 0 {
		return errNoSkus
	}

	// 3. Convert response
	var outputItems []ItemOutput
	for _, inputItem := range inputItems {
		existed := false
		for _, product := range products {
			if inputItem.SellerSku == product.SellerSku && inputItem.UomName == product.UomName {
				itemOutput := ItemOutput{product.Sku, inputItem.Quantity}
				outputItems = append(outputItems, itemOutput)
				existed = true
				break
			}
		}
		if !existed {
			msg := fmt.Sprintf("không tìm thấy thông tin sku với sellerSku=%s, uomName=%s", inputItem.SellerSku, inputItem.UomName)
			return customFunc.FuncResult{ErrorMessage: msg}
		}
	}

	return customFunc.FuncResult{Result: outputItems}
}

func callApiGetSkus(subItems []ItemInput) ([]Product, error) {
	// 1. Convert input to param
	sellerSkus := converter.Map(subItems, func(i ItemInput) string { return i.SellerSku })
	sellerSkusStr := strings.Join(sellerSkus[:], ",")
	// todo batch 50

	// 2. Prepare call api
	httpClient := initHttpClient()
	reqHeader := map[string]string{"Content-Type": "application/json"}
	reqParams := map[string]string{"sellerSkus": sellerSkusStr}

	// 3. Call api
	httpStatus, resBody, err := utils.SendHTTPRequest[any, GetSkuResponse](httpClient, http.MethodGet, urlApiGetSkus, reqHeader, reqParams, nil)
	if err != nil || httpStatus != http.StatusOK {
		logger.Errorf("failed to call %v, got error=%v, httpStatus=%d, resBody=%+v", urlApiGetSkus, err, httpStatus, resBody)
		return nil, err
	}

	// 4. Return data
	return resBody.Data.Products, nil
}

// initHttpClient...
func initHttpClient() *http.Client {
	transportCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   20 * time.Second,
		Transport: transportCfg,
	}
	return client
}
