package configloader

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
)

// ConfigLoader ........................................................................................................
type ConfigLoader interface {
	// Load ...
	Load(file fileprocessing.ProcessingFile) (ConfigMappingMD, error)
}
