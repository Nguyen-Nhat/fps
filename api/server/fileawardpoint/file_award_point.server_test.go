package fileawardpoint

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAPIListFile_validateAndSetDataValue_Size_bigger_than_200_with_valid_page(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/lfp/fileAwardPoint/getList?page=1&size=300", nil)

	_, err := validateAndSetDataValue(req)

	assert.Equal(t, fmt.Errorf("size out of range"), err)
}

func TestAPIListFile_validateAndSetDataValue_Size_bigger_than_200_without_page(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/lfp/fileAwardPoint/getList?size=300", nil)

	_, err := validateAndSetDataValue(req)

	assert.Equal(t, fmt.Errorf("size out of range"), err)
}

func TestAPIListFile_validateAndSetDataValue_Page_bigger_than_1000_with_valid_size(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/lfp/fileAwardPoint/getList?page=1001&size=1", nil)

	_, err := validateAndSetDataValue(req)

	assert.Equal(t, fmt.Errorf("page out of range"), err)
}

func TestAPIListFile_validateAndSetDataValue_Page_bigger_than_1000_without_size(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:10080/lfp/fileAwardPoint/getList?page=1001", nil)

	_, err := validateAndSetDataValue(req)

	assert.Equal(t, fmt.Errorf("page out of range"), err)
}
