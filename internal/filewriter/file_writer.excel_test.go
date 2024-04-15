package filewriter

import (
	"testing"
)

func Test_excelFileWriter_OutputFileContentType(t *testing.T) {
	c := &excelFileWriter{outputFileContentType: "abcd"}
	if got := c.OutputFileContentType(); got != c.outputFileContentType {
		t.Errorf("OutputFileContentType() = %v, want %v", got, c.outputFileContentType)
	}
}
