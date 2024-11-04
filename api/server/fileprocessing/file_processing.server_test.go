package fileprocessing

import (
	"fmt"
	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAPIListFile_validateAndSetDataValue_Size_bigger_than_200_with_valid_page(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/v1/getListProcessFiles?clientId=1&page=1&size=300", nil)
	data := &GetFileProcessHistoryRequest{}
	err := bindAndValidateRequestParams(req, data)

	expectErr := error2.ErrInvalidRequestWithError(fmt.Errorf("request field is out of range: size"))
	assert.Equal(t, expectErr, err)
}

func TestAPIListFile_validateAndSetDataValue_Size_bigger_than_200_without_page(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/v1/getListProcessFiles?clientId=1&size=300", nil)
	data := &GetFileProcessHistoryRequest{}
	err := bindAndValidateRequestParams(req, data)

	expectErr := error2.ErrInvalidRequestWithError(fmt.Errorf("request field is out of range: size"))
	assert.Equal(t, expectErr, err)
}

func TestAPIListFile_validateAndSetDataValue_Page_bigger_than_1000_with_valid_size(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/v1/getListProcessFiles?clientId=1&page=1001&size=1", nil)
	data := &GetFileProcessHistoryRequest{}
	err := bindAndValidateRequestParams(req, data)

	expectErr := error2.ErrInvalidRequestWithError(fmt.Errorf("request field is out of range: page"))
	assert.Equal(t, expectErr, err)
}

func TestAPIListFile_validateAndSetDataValue_Page_bigger_than_1000_without_size(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/v1/getListProcessFiles?clientId=1&page=1001", nil)
	data := &GetFileProcessHistoryRequest{}
	err := bindAndValidateRequestParams(req, data)

	expectErr := error2.ErrInvalidRequestWithError(fmt.Errorf("request field is out of range: page"))
	assert.Equal(t, expectErr, err)
}

func TestAPIListFile_validateAndSetDataValue_Page_not_include_clientId(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/v1/getListProcessFiles?page=1&size=1", nil)
	data := &GetFileProcessHistoryRequest{}
	err := bindAndValidateRequestParams(req, data)

	expectErr := error2.ErrInvalidRequestWithError(fmt.Errorf("missing required param: clientId or clientIds"))
	assert.Equal(t, expectErr, err)
}
