package customFunc

import (
	"reflect"
	"testing"
)

func Test_reUploadFile(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want FuncResult
	}{
		{"test reUploadFile", "https://abc.com/image.jpg", FuncResult{"https://abc.com/image.jpgabc", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reUploadFile(tt.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reUploadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
