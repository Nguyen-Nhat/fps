package loyalty

import (
	"fmt"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"reflect"
	"testing"
)

// For temporary check ... removeThisPrefixForRunningThisTest_
func TestClient_GrantPoint(t *testing.T) {
	cfg := config.LoyaltyConfig{
		Endpoint: "http://127.0.0.1:8080",
		APIKey:   "nguyenducquocdai",
		Paths: config.LoyaltyPaths{
			TxnGrant: "/api/v4/transaction/grantPoint",
		},
	}
	loyaltyClient := NewClient(cfg)

	req := GrantPointRequest{
		MerchantID: "247954728103710720",
		RefId:      "2379842379474",
		Phone:      "0393227489",
		Point:      1,
		TxnDesc:    "",
	}
	res, err := loyaltyClient.GrantPoint(req)
	if err != nil {
		logger.Errorf("error %v", err)
	}
	if res.IsFailed() {
		logger.Errorf("error %v", res.ToStringCodeMessage())
	}
	if res.IsSuccess() {
		fmt.Printf("Success %v", res.Data)
	}
}

// For temporary check ... removeThisPrefixForRunningThisTest_
func removeThisPrefixForRunningThisTest_TestClient_GetTransactionByID(t *testing.T) {
	cfg := config.LoyaltyConfig{
		Endpoint: "http://127.0.0.1:8080",
		APIKey:   "nguyenducquocdai",
		Paths: config.LoyaltyPaths{
			TxnGetList: "/api/v4/transaction/getListTransaction",
		},
	}
	loyaltyClient := NewClient(cfg)

	res, err := loyaltyClient.GetTransactionByID(378112751274299392)
	if err != nil {
		logger.Errorf("error %v", err)
	}
	if res.IsFailed() {
		logger.Errorf("error %v", res.ToStringCodeMessage())
	}
	if res.IsSuccess() {
		fmt.Printf("Success %v", res.Data)
	}
}

func Test_defaultHeader(t *testing.T) {
	type args struct {
		cfg config.LoyaltyConfig
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Header must contains application_json type and correct api key",
			args: args{cfg: config.LoyaltyConfig{APIKey: "my_api_key"}},
			want: map[string]string{"Content-Type": "application/json", "X-API-KEY": "my_api_key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultHeader(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
