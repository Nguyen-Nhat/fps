package errorz

import "fmt"

func SheetNotFound(sheetName string) error {
	return fmt.Errorf("sheet %s not found", sheetName)
}
