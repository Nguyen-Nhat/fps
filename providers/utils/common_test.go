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
			got, err := GenerateRandomBytes(tt.args.byteLength)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractFileName(tt.args.filePath)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filename = %#v, want %#v", got, tt.want)
			}
		})
	}
}
