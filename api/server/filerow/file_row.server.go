package filerow

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	commonError "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	res "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	IServer interface {
		GetListFileRowAPI() func(http.ResponseWriter, *http.Request)
		GetListFileRow(ctx context.Context, fileID int, req *GetListFileRowData) (*res.BaseResponse[GetListFileRowResponse], error)
	}

	// Server ...
	Server struct {
		service fileprocessingrow.Service
	}
)

// *Server implements IServer
var _ IServer = &Server{}

// InitFileRowServer ...
func InitFileRowServer(db *sql.DB) *Server {
	repo := fileprocessingrow.NewRepo(db)
	service := fileprocessingrow.NewService(repo)
	return &Server{
		service: service,
	}
}

func (s *Server) GetListFileRowAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get file ID
		fileIDStr := chi.URLParam(r, "fileID")
		logger.Infof("GetDetailClient: fileID get from url %s", fileIDStr)
		fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
		if err != nil {
			logger.Errorf("===== GetListFileRowAPI: fileID is not int", err.Error())
			_ = render.Render(w, r, commonError.ErrRenderInvalidRequest(err))
			return
		}

		// 1.Bind data & validate input
		data, err := bindAndValidateRequestParams(r)
		if err != nil {
			logger.Errorf("===== GetListFileRowAPI: Bind data and validate input error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ErrRenderInvalidRequest(err))
			return
		}

		// 2. Handle request
		resp, err := s.GetListFileRow(r.Context(), int(fileID), data)
		if err != nil {
			logger.Errorf("===== GetListFileRowAPI handler error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}

		// 3. Render response
		logger.Infof("===== GetListFileRowAPI response: size=%v\n", len(resp.Data.Rows))
		err = render.Render(w, r, resp)
		if err != nil {
			logger.Errorf("===== GetListFileRowAPI render response error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}
	}
}

func (s *Server) GetListFileRow(ctx context.Context, fileID int, req *GetListFileRowData) (*res.BaseResponse[GetListFileRowResponse], error) {
	// 1. Map the request to internal DTO
	input := fileprocessingrow.GetListFileRowsRequest{
		PageRequest: req.PageRequest,
		// ...
	}

	// 2. Handle request
	fps, pagination, err := s.service.GetListFileRowsByFileID(ctx, fileID, input)
	if err != nil {
		return nil, err
	}

	// 3. Return
	resp := &GetListFileRowResponse{
		Rows:       fps,
		Pagination: pagination,
	}
	resp2 := res.ToResponse(resp)
	return resp2, nil
}
