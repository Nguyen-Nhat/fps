package fpsclient

import "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/request"

// GetListClientDTO ...
type GetListClientDTO struct {
	request.PageRequest
	Name string
}
