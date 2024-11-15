package errorz

import "fmt"

func SheetNotFound(sheetName string) error {
	return fmt.Errorf("sheet %s not found", sheetName)
}

func ErrCantExtractFileExtensionFromURL(fileURL string) error {
	return fmt.Errorf("cannot extract file extension from file URL %v", fileURL)
}
