package fpsclient

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/render"

	commonError "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	res "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fpsclient"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

type (
	IServer interface {
		GetListClientAPI() func(http.ResponseWriter, *http.Request)
		GetListClient(ctx context.Context, req *GetListClientData) (*res.BaseResponse[GetListClientResponse], error)
	}

	// Server ...
	Server struct {
		service fpsclient.Service
	}
)

// *Server implements IServer
var _ IServer = &Server{}

// InitClientServer ...
func InitClientServer(db *sql.DB) *Server {
	repo := fpsclient.NewRepo(db)
	service := fpsclient.NewService(repo)
	return &Server{
		service: service,
	}
}

func (s *Server) GetListClientAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.Bind data & validate input
		data, err := bindAndValidateRequestParams(r)
		if err != nil {
			logger.Errorf("===== GetListClientAPI: Bind data and validate input error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ErrRenderInvalidRequest(err))
			return
		}

		// 2. Handle request
		resp, err := s.GetListClient(r.Context(), data)
		if err != nil {
			logger.Errorf("===== GetListClientAPI handler error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}

		// 3. Render response
		logger.Infof("===== GetListClientAPI response: \n%v\n", utils.JsonString(resp))
		err = render.Render(w, r, resp)
		if err != nil {
			logger.Errorf("===== GetListClientAPI render response error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}
	}
}

func (s *Server) GetListClient(ctx context.Context, req *GetListClientData) (*res.BaseResponse[GetListClientResponse], error) {
	// 1. Map the request to internal DTO
	input := fpsclient.GetListClientDTO{
		Name:        req.Name,
		PageRequest: req.PageRequest,
	}

	// 2. Handle request
	fps, pagination, err := s.service.GetListClients(ctx, input)
	if err != nil {
		return nil, err
	}
	// 3. Return
	resp := toGetListClientResponse(fps, pagination)
	resp2 := res.ToResponse(resp)
	return resp2, nil
}
