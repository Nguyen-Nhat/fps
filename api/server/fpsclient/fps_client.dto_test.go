package fpsclient

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/request"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

func Test_bindAndValidateRequestParams(t *testing.T) {
	randomStr255 := utils.RandStringBytes(255)
	randomStr256 := utils.RandStringBytes(256)

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *GetListClientData
		wantErr bool
	}{
		{"test query no value -> return data with default value", args{toRequest("", 0, 0)},
			&GetListClientData{request.PageRequest{Page: 1, PageSize: 10}, ""}, false},
		{"test query with name -> return data with name", args{toRequest("abc", 0, 0)},
			&GetListClientData{request.PageRequest{Page: 1, PageSize: 10}, "abc"}, false},
		{"test query with name, page, pageSize -> return data with name", args{toRequest("abc", 2, 20)},
			&GetListClientData{request.PageRequest{Page: 2, PageSize: 20}, "abc"}, false},
		{"test query with name is 255 -> return error", args{toRequest(randomStr255, 2, 123)},
			&GetListClientData{request.PageRequest{Page: 2, PageSize: 123}, randomStr255}, false},

		{"test query with negative page -> return error", args{toRequest("", -1, 20)},
			nil, true},
		{"test query with negative pageSize -> return error", args{toRequest("", 2, -123)},
			nil, true},
		{"test query with name is 256 -> return error", args{toRequest(randomStr256, 2, 12)},
			nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bindAndValidateRequestParams(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("bindAndValidateRequestParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bindAndValidateRequestParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func toRequest(name string, page int, pageSize int) *http.Request {
	domain := "http://localhost:1000/abc.com?name=%v&page=%v&pageSize=%v"
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(domain, name, page, pageSize), nil)
	return req
}
