package response

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

type BaseResponse[D any] struct {
	Error   codes.Code `json:"code"`
	Message string     `json:"message"`
	Data    *D         `json:"data"`
}

func ToResponse[D any](data *D) *BaseResponse[D] {
	return &BaseResponse[D]{Data: data}
}

func (r *BaseResponse[D]) Render(_ http.ResponseWriter, _ *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	r.Error = codes.OK
	r.Message = "Successfully"
	return nil
}
