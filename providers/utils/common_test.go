package utils

import (
	"errors"
	"reflect"
	"testing"
)

func TestHiddenString(t *testing.T) {
	type args struct {
		input            string
		numberOfTailChar int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "input is empty",
			args: args{"", 3},
			want: "",
		},
		{
			name: "numberOfTailChar is zero",
			args: args{"abc", 0},
			want: "***",
		},
		{
			name: "input shorter than tail",
			args: args{"ab", 3},
			want: "ab",
		},
		{
			name: "input longer than tail",
			args: args{"abcd", 3},
			want: "***bcd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HiddenString(tt.args.input, tt.args.numberOfTailChar); got != tt.want {
				t.Errorf("HiddenString() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGenerateRandomString(t *testing.T) {
	type args struct {
		byteLength int
	}
	tests := []struct {
		name         string
		args         args
		err          error
		wantedLength int
	}{
		{
			name:         "generate not empty string if byte length != 0",
			args:         args{10},
			wantedLength: 14,
			err:          nil,
		},
		{
			name:         "generate empty string if byte length = 0",
			args:         args{0},
			wantedLength: 0,
			err:          nil,
		},
		{
			name:         "err is nil",
			args:         args{10},
			wantedLength: 14,
			err:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomString(tt.args.byteLength)
			if !errors.Is(err, tt.err) {
				t.Errorf("HiddenString() = %v, want err %v", err, tt.err)
			}
			if len(got) != tt.wantedLength {
				t.Errorf("Length = %v, want %v", len(got), tt.wantedLength)
			}
		})
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	type args struct {
		byteLength int
	}
	tests := []struct {
		name         string
		args         args
		err          error
		wantedLength int
	}{
		{
			name:         "generate empty byte",
			args:         args{0},
			wantedLength: 0,
			err:          nil,
		},
		{
			name:         "generate byte length 1",
			args:         args{1},
			wantedLength: 1,
			err:          nil,
		},
		{
			name:         "generate byte length 2",
			args:         args{2},
			wantedLength: 2,
			err:          nil,
		},
		{
			name:         "generate byte length 5",
			args:         args{5},
			wantedLength: 5,
			err:          nil,
		},
		{
			name:         "generate byte length 999999",
			args:         args{999999},
			wantedLength: 999999,
			err:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateRandomBytes(tt.args.byteLength)
			if !errors.Is(err, tt.err) {
				t.Errorf("HiddenString() = %v, want err %v", err, tt.err)
			}
			if len(got) != tt.wantedLength {
				t.Errorf("Length = %v, want %v", len(got), tt.wantedLength)
			}
		})
	}
}

func TestExtractFileName(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want FileName
		skip bool
	}{
		{
			name: "get file name from path",
			args: args{"abc/xyz.doc"},
			want: FileName{
				FullName:  "xyz.doc",
				Name:      "xyz",
				Extension: "doc",
			},
		},
		{
			name: "get file name from file name",
			args: args{"xyz.doc"},
			want: FileName{
				FullName:  "xyz.doc",
				Name:      "xyz",
				Extension: "doc",
			},
		},
		{
			name: "not panic when meet invalid path",
			args: args{"a/b/c/d"},
			want: FileName{
				FullName:  "unknown",
				Name:      "unknown",
				Extension: "unknown",
			},
		},
		{
			name: "get file name from multiple dot fileName",
			args: args{"abc.d.e.doc"},
			want: FileName{
				FullName:  "abc.d.e.doc",
				Name:      "abc.d.e",
				Extension: "doc",
			},
		},
		{
			name: "get file name from multiple dot path",
			args: args{"folder/folder2/abc.d.e.doc"},
			want: FileName{
				FullName:  "abc.d.e.doc",
				Name:      "abc.d.e",
				Extension: "doc",
			},
		},
		{
			name: "get file name from path",
			args: args{"https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/19/4d9ee61c-b725-47ef-b652-4e7bbf1c57fb/VietMeta%20-%20nap%20diem%20KH%202022-10-19.xlsx"},
			want: FileName{
				FullName:  "VietMeta - nap diem KH 2022-10-19.xlsx",
				Name:      "VietMeta - nap diem KH 2022-10-19",
				Extension: "xlsx",
			},
		},
		{
			name: "get file name from path",
			args: args{"https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/19/7f60d01e-1e0a-47ab-a7fc-4c8e4816f7a0/VietMeta%20-%20nap%20diem%20KH%202022-10-19%20%28part%204%29.xlsx"},
			want: FileName{
				FullName:  "VietMeta - nap diem KH 2022-10-19 (part 4).xlsx",
				Name:      "VietMeta - nap diem KH 2022-10-19 (part 4)",
				Extension: "xlsx",
			},
		},
		{
			name: "get file name from name with space",
			args: args{"VietMeta - nap diem KH 2022-10-19 (part 4).xlsx"},
			want: FileName{
				FullName:  "VietMeta - nap diem KH 2022-10-19 (part 4).xlsx",
				Name:      "VietMeta - nap diem KH 2022-10-19 (part 4)",
				Extension: "xlsx",
			},
		},
		{
			name: "get file name from name with space",
			args: args{"Vietmeta Loyalty - Nạp điểm.xlsx"},
			want: FileName{
				FullName:  "Vietmeta Loyalty - Nạp điểm.xlsx",
				Name:      "Vietmeta Loyalty - Nạp điểm",
				Extension: "xlsx",
			},
		},
		{
			name: "get file name from private url",
			args: args{"https://files.dev.tekoapis.net/files/a4c7aca4-a6bb-4ce6-b4be-c5c411da5fca"},
			want: FileName{
				FullName:  "sc_export_specialist_05032024205328.xlsx",
				Name:      "sc_export_specialist_05032024205328",
				Extension: "xlsx",
			},
			skip: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}
			got := ExtractFileName(tt.args.filePath)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filename = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestRandStringBytes(t *testing.T) {
	type args struct {
		numberChars int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test random string with length = 0",
			args: args{numberChars: 0},
			want: 0,
		},
		{
			name: "test random string with length = 1",
			args: args{numberChars: 1},
			want: 1,
		},
		{
			name: "test random string with length = 100",
			args: args{numberChars: 100},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandStringBytes(tt.args.numberChars); len(got) != tt.want {
				t.Errorf("RandStringBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileName_FullNameWithSuffix(t *testing.T) {
	type fields struct {
		FullName  string
		Name      string
		Extension string
	}
	tests := []struct {
		name   string
		fields fields
		suffix string
		want   string
	}{
		{"test FullNameWithSuffix case normal", fields{"full_name", "name", "abc"}, "_result", "name_result.abc"},
		{"test FullNameWithSuffix case name empty", fields{"", "", "abc"}, "_result", "_result.abc"},
		{"test FullNameWithSuffix case suffix empty", fields{"full_name", "name", "abc"}, "", "name.abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileName{
				FullName:  tt.fields.FullName,
				Name:      tt.fields.Name,
				Extension: tt.fields.Extension,
			}
			if got := f.FullNameWithSuffix(tt.suffix); got != tt.want {
				t.Errorf("FullNameWithSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TrimSpaceAndToLower(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"test TrimSpaceAndToLower case empty", "", ""},
		{"test TrimSpaceAndToLower case normal", "  This is    a string with space and UPPER and lower     ", "this is    a string with space and upper and lower"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimSpaceAndToLower(tt.input); got != tt.want {
				t.Errorf("TrimSpaceAndToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualsIgnoreCase(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"", ""}, true},
		{"", args{"  abc  ", "abc   "}, true},
		{"", args{"  abBDc  ", "abbDc   "}, true},
		{"", args{"  this Is A test  ", "               this iS a TEST   "}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualsIgnoreCase(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("EqualsIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResultFileName(t *testing.T) {
	type args struct {
		fileName string
		fileExt  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"file name with mime type", args{"abc.xlsx", ""}, "abc_result.xlsx"},
		{"file name without mime type", args{"abc", ""}, "abc_result"},
		{"file name with mime type config", args{"abc", "xlsx"}, "abc_result.xlsx"},
		{"file name empty", args{"", "xlsx"}, ""},
		{"file name with other mime type", args{"abc.xlsx", "csv"}, "abc_result.csv"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetResultFileName(tt.args.fileName, tt.args.fileExt); got != tt.want {
				t.Errorf("GetResultFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IndexOf(t *testing.T) {
	type args[T comparable] struct {
		arr []T
		val T
	}
	type testCase[T comparable] struct {
		name  string
		args  args[T]
		want  int
		want1 bool
	}
	tests := []testCase[int]{
		{"test IndexOf: array is nil -> false", args[int]{nil, 1}, -1, false},
		{"test IndexOf: array is empty -> false", args[int]{[]int{}, 1}, -1, false},
		{"test IndexOf: array has data, but not found -> false", args[int]{[]int{1, 2, 3, 4, 5, 6, 7}, 10}, -1, false},
		{"test IndexOf: array has data, and value is existed -> false", args[int]{[]int{1, 2, 3, 4, 5, 6, 7}, 3}, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := IndexOf(tt.args.arr, tt.args.val)
			if got != tt.want {
				t.Errorf("IndexOf() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IndexOf() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
