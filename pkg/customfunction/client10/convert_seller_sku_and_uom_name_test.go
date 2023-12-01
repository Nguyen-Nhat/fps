package funcClient10

import (
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

func Test_GetSkus(t *testing.T) {
	res := ConvertSellerSkus(jsonStr, "12")
	logger.Infof("Result = %+v", res)
}

const jsonStr = `
[
  {
    "sellerSku": "1703401",
    "uomName": "Cái",
    "requestQty": 0.8
  }
]
`
