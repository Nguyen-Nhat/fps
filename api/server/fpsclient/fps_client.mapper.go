package fpsclient

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fpsclient"
)

func toGetListClientResponse(clients []*fpsclient.Client, pagination response.PaginationNew) *GetListClientResponse {
	result := make([]ClientDTO, 0)
	for _, v := range clients {
		result = append(result, toClientDTO(v))
	}

	return &GetListClientResponse{
		Clients:    result,
		Pagination: pagination,
	}

}

func toClientDTO(client *fpsclient.Client) ClientDTO {
	return ClientDTO{
		ID:            client.ID,
		Name:          client.Name,
		Description:   client.Description,
		SampleFileURL: client.SampleFileURL,
		CreatedAt:     client.CreatedAt.UnixMilli(),
		CreatedBy:     client.CreatedBy,
	}
}
