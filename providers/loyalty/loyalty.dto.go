package loyalty

import (
	"fmt"
	"time"
)

// Common --------------------------------------------------------------------------------------------------------------

const (
	responseSuccessCode      = "00"
	responseNotEnoughBalance = "4280"

	txnSuccessStatus = "4"
)

type BaseLsResponse[D any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    *D     `json:"data"`
}

func (r *BaseLsResponse[any]) IsSuccess() bool {
	return r.Code == responseSuccessCode
}

func (r *BaseLsResponse[any]) IsNotEnoughBalance() bool {
	return r.Code == responseNotEnoughBalance
}

func (r *BaseLsResponse[any]) IsFailed() bool {
	return !r.IsSuccess()
}

func (r *BaseLsResponse[any]) ToStringCodeMessage() string {
	return fmt.Sprintf("%v - %v", r.Code, r.Message)
}

type PaginationData struct {
	CurrentPage int32 `json:"currentPage"`
	PageSize    int32 `json:"pageSize"`
	TotalItems  int64 `json:"totalItems"`
	TotalPage   int32 `json:"totalPage"`
}

// Grant Point Request & Response --------------------------------------------------------------------------------------

// Request ------

type GrantPointRequest struct {
	MerchantID string
	RefId      string
	Phone      string
	Point      int32
	TxnDesc    string
}

type grantPointRequest struct {
	MerchantID string  `json:"merchantId"`
	GrantType  int32   `json:"grantType"`
	Amount     float64 `json:"amount"`
	Point      int32   `json:"point"`
	Phone      string  `json:"phone"`
	OrderCode  string  `json:"orderCode"`
	RefID      string  `json:"refId"`
	RefTime    int64   `json:"refTime"`
	TxnDesc    string  `json:"txnDesc"`
}

func (r *grantPointRequest) mapByRequest(req GrantPointRequest) {
	r.MerchantID = req.MerchantID
	r.GrantType = 2
	r.Point = req.Point
	r.Phone = req.Phone
	r.RefID = req.RefId
	r.RefTime = time.Now().Truncate(time.Second).UnixMilli()
	r.TxnDesc = req.TxnDesc
}

// Response ---------

type GrantPointResponse struct {
	Transaction GrantTxnResponse `json:"transaction"`
}

type GrantTxnResponse struct {
	TxnID string `json:"txnId"`
	Point string `json:"point"`
}

// Get List transaction Request & Response -----------------------------------------------------------------------------

type getListTxnRequest struct {
	Ids       *[]int64 `json:"ids"`
	NetworkID *int32   `json:"networkId"`
	Phone     *string  `json:"phone"`
	FromTime  *int64   `json:"fromTime"`
	ToTime    *int64   `json:"toTime"`
	Status    *int32   `json:"status"`
	TxnTypes  []string `json:"txnTypes"`
	Page      *int32   `json:"page"`
	Size      *int32   `json:"size"`
}

func (r *getListTxnRequest) mapByID(id int64) {
	r.Ids = &[]int64{id}
}

type GetListTxnResponse struct {
	Pagination   PaginationData     `json:"pagination"`
	Transactions []TransactionsData `json:"transactions"`
}

type TransactionsData struct {
	TxnID          string `json:"txnId"`
	RefID          string `json:"refId"`
	OrderCode      string `json:"orderCode"`
	MerchantID     string `json:"merchantId"`
	MerchantName   string `json:"merchantName"`
	TxnType        string `json:"txnType"`
	TxnTypeText    string `json:"txnTypeText"`
	Point          string `json:"point"`
	TxnDesc        string `json:"txnDesc"`
	AccountingType string `json:"accountingType"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"`
	RefTime        string `json:"refTime"`
	CreatedBy      string `json:"createdBy"`
}

func (t *TransactionsData) IsSuccess() bool {
	return t.Status == txnSuccessStatus
}
