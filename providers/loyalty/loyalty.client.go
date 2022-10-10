package loyalty

import (
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"net/http"
	"time"
)

type (
	// IClient ...
	IClient interface {
		GrantPoint(GrantPointRequest) (BaseLsResponse[GrantPointResponse], error)
		GetTransactionByID(int64) (BaseLsResponse[GetListTxnResponse], error)
	}

	// Client ....
	Client struct {
		conf   config.LoyaltyConfig
		client *http.Client
	}
)

var _ IClient = &Client{}

// NewClient ...
func NewClient(conf config.LoyaltyConfig) *Client {
	if len(conf.Endpoint) == 0 {
		panic("===== Loyalty Endpoint Must Not Empty")
	}

	return &Client{
		conf: conf,
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *Client) GrantPoint(req GrantPointRequest) (BaseLsResponse[GrantPointResponse], error) {
	// 1. Build url & header
	url := c.conf.Endpoint + c.conf.Paths.TxnGrant
	header := defaultHeader(c.conf)

	// 2. Build request body
	loyRequest := grantPointRequest{}
	loyRequest.mapByRequest(req)

	// 3. Send http request
	httpResp, err := utils.SendHTTPRequest[grantPointRequest, BaseLsResponse[GrantPointResponse]](c.client, http.MethodPost, url, header, &loyRequest)
	if err != nil {
		logger.Errorf("===== loyalty.GrantPoint: Call Loyalty Error: %s", err.Error())
		return BaseLsResponse[GrantPointResponse]{}, err
	}

	// 5. Response
	return *httpResp, nil
}

func (c *Client) GetTransactionByID(id int64) (BaseLsResponse[GetListTxnResponse], error) {
	// 1. Build url & header
	url := c.conf.Endpoint + c.conf.Paths.TxnGetList
	header := defaultHeader(c.conf)

	// 2. Build request body
	request := getListTxnRequest{}
	request.mapByID(id)

	// 3. Send http request
	httpResp, err := utils.SendHTTPRequest[getListTxnRequest, BaseLsResponse[GetListTxnResponse]](c.client, http.MethodPost, url, header, &request)
	if err != nil {
		logger.Errorf("===== loyalty.GetTransactionByID: Call Loyalty Error: %s", err.Error())
		return BaseLsResponse[GetListTxnResponse]{}, err
	}

	// 5. Response
	return *httpResp, nil
}

// private method ------------------------------------------------------------------------------------------------------

func defaultHeader(cfg config.LoyaltyConfig) map[string]string {
	apiKey := cfg.APIKey

	subApiKey := utils.HiddenString(apiKey, 10)
	logger.Infof("loyalty.header: api_key = %v", subApiKey)

	return map[string]string{
		"Content-Type": "application/json",
		"X-API-KEY":    apiKey,
	}
}
