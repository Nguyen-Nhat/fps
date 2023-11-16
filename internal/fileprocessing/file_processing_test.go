package fileprocessing

import (
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

func TestProcessingFile_IsInitStatus(t *testing.T) {
	tests := []struct {
		name           string
		ProcessingFile ent.ProcessingFile
		want           bool
	}{
		{"test IsInitStatus case Init", ent.ProcessingFile{Status: StatusInit}, true},
		{"test IsInitStatus case Processing", ent.ProcessingFile{Status: StatusProcessing}, false},
		{"test IsInitStatus case Failed", ent.ProcessingFile{Status: StatusFailed}, false},
		{"test IsInitStatus case Finished", ent.ProcessingFile{Status: StatusFinished}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := &ProcessingFile{ProcessingFile: tt.ProcessingFile}
			if got := pf.IsInitStatus(); got != tt.want {
				t.Errorf("IsInitStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessingFile_IsProcessingStatus(t *testing.T) {
	tests := []struct {
		name           string
		ProcessingFile ent.ProcessingFile
		want           bool
	}{
		{"test IsInitStatus case Init", ent.ProcessingFile{Status: StatusInit}, false},
		{"test IsInitStatus case Processing", ent.ProcessingFile{Status: StatusProcessing}, true},
		{"test IsInitStatus case Failed", ent.ProcessingFile{Status: StatusFailed}, false},
		{"test IsInitStatus case Finished", ent.ProcessingFile{Status: StatusFinished}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := &ProcessingFile{
				ProcessingFile: tt.ProcessingFile,
			}
			if got := pf.IsProcessingStatus(); got != tt.want {
				t.Errorf("IsProcessingStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	if got := Name(); got != "ProcessingFile" {
		t.Errorf("Name() = %v, want %v", got, "ProcessingFile")
	}
}
