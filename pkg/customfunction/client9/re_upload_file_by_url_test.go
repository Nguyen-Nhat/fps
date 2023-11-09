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
//	//fileUrl := "https://drive.google.com/file/d/1SNtCNwVmN3Hvf5tbOnmvEIfqdB9utVeR/view?usp=drive_link"
//	result := ReUploadFile(fileUrl)
//	logger.Infof("%v", result)
//}

func Test_extractGoogleDriveLink(t *testing.T) {
	tests := []struct {
		name         string
		fileUrl      string
		wantFixedUrl string
		wantFileName string
		wantErr      bool
	}{
		// case wrong
		{"test extractGoogleDriveLink case url is only domain drive.google.com", googleDriveUrl, "", "", true},
		{"test extractGoogleDriveLink case url is wrong 1", "https://drive.google.com", "", "", true},
		{"test extractGoogleDriveLink case url is wrong 2", "https://drive.google.com/", "", "", true},

		// check with file/u
		{"test extractGoogleDriveLink case url not start with file/u", "https://drive.google.com/file/a/1/2/abcde", "", "", true},
		{"test extractGoogleDriveLink case url starts with file/u, but no fileID 1", "https://drive.google.com/file/u/", "", "", true},
		{"test extractGoogleDriveLink case url starts with file/u, but no fileID 2", "https://drive.google.com/file/u/1", "", "", true},
		{"test extractGoogleDriveLink case url starts with file/u, but no fileID 3", "https://drive.google.com/file/u/1/d/", "", "", true},

		{"test extractGoogleDriveLink case url starts with file/u, and has fileID 1", "https://drive.google.com/file/u/1/d/abcde",
			"https://drive.google.com/uc?id=abcde&export=download", "abcde.jpg", false},
		{"test extractGoogleDriveLink case url starts with file/u, and has fileID 2", "https://drive.google.com/file/u/1/m/abcde",
			"https://drive.google.com/uc?id=abcde&export=download", "abcde.jpg", false},
		{"test extractGoogleDriveLink case url starts with file/u, and has fileID 3", "https://drive.google.com/file/u/1/m/abcde/mnp",
			"https://drive.google.com/uc?id=abcde&export=download", "abcde.jpg", false},

		// check with file/d
		{"test extractGoogleDriveLink case url not start with file/d", "https://drive.google.com/file/a/abcde", "", "", true},
		{"test extractGoogleDriveLink case url starts with file/d, but no fileID 1", "https://drive.google.com/file/d/", "", "", true},
		{"test extractGoogleDriveLink case url starts with file/d, and has fileID 1", "https://drive.google.com/file/d/abcde",
			"https://drive.google.com/uc?id=abcde&export=download", "abcde.jpg", false},
		{"test extractGoogleDriveLink case url starts with file/d, and has fileID 2", "https://drive.google.com/file/d/abcde/mnpa",
			"https://drive.google.com/uc?id=abcde&export=download", "abcde.jpg", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := extractGoogleDriveLink(tt.fileUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractGoogleDriveLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantFixedUrl {
				t.Errorf("extractGoogleDriveLink() got = %v, want %v", got, tt.wantFixedUrl)
			}
			if got1 != tt.wantFileName {
				t.Errorf("extractGoogleDriveLink() got1 = %v, want %v", got1, tt.wantFileName)
			}
		})
	}
}
