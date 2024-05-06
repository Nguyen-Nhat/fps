package xls

import (
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"go.tekoapis.com/tekone/library/test/monkey"
)

type xlsTestSuite struct {
	suite.Suite
}

func TestXls(t *testing.T) {
	ts := &xlsTestSuite{}

	monkey.PatchInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(_ *http.Client, url string) (*http.Response, error) {
		file, err := os.Open(url)
		return &http.Response{
			Body: file,
		}, err
	})

	suite.Run(t, ts)
}

func (ts *xlsTestSuite) assert(url, sheetName, wantFile string) {
	data, err := LoadXlsByUrl(url, sheetName)
	if err != nil {
		goldie.New(ts.T()).AssertJson(ts.T(), wantFile, err.Error())
		return
	}

	assert.Nil(ts.T(), err)
	goldie.New(ts.T()).AssertJson(ts.T(), wantFile, data)
}

func (ts *xlsTestSuite) Test200_EmptySheetName_GetDefaultFirstSheet() {
	ts.assert("./testdata/happy_case.xls", "", "happy_case_empty_sheet_name_get_default_first_sheet")
}

func (ts *xlsTestSuite) Test200_PassSheetName_ThenReturnSuccess() {
	ts.assert("./testdata/happy_case.xls", "data", "happy_case_pass_sheet_name_then_return_success")
}
