package funcClient9

import (
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
)

func Test_reUploadFile(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want customFunc.FuncResult
	}{
		{"test reUploadFile with empty url", "", customFunc.FuncResult{Result: nil}},
		{"test reUploadFile", "https://abc.com/image.jpg", errDefault},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReUploadFile(tt.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reUploadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_uploadFile(t *testing.T) {
//	fileUrl := "https://cdn.nguyenkimmall.com/images/thumbnails/600/336/detailed/873/10055363-may-hut-bui-khong-day-tefal-x-nano-ty1129wo-2.jpg"
//	//fileUrl := "https://drive.google.com/file/u/1/d/1SNtCNwVmN3Hvf5tbOnmvEIfqdB9utVeR/view?usp=drive_link"
//	result := ReUploadFile(fileUrl)
//	logger.Infof("%v", result)
//}
