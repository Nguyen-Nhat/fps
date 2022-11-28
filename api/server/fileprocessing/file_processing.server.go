package fileprocessing

import (
	"context"
	"database/sql"
	"net/http"

	"fmt"
	commonError "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	res "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/go-chi/render"
	"strconv"
)

type (
	IServer interface {
		GetFileProcessHistoryAPI() func(http.ResponseWriter, *http.Request)
		CreateProcessByFileAPI() func(http.ResponseWriter, *http.Request)
	}

	// Server ...
	Server struct {
		service *fileprocessing.ServiceImpl
	}
)

func (s *Server) GetFileProcessHistoryAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := validateParameterAndSetDataValue(r)
		if err != nil {
			render.Render(w, r, error2.ErrInvalidRequest(err))
			return
		}

		resp, err := s.GetFileProcessHistory(r.Context(), data)
		if err != nil {
			render.Render(w, r, error2.ErrInternal(err))
			return
		}

		render.Render(w, r, resp)
	}
}

func (s *Server) GetFileProcessHistory(ctx context.Context, req *fileprocessing.GetFileProcessHistoryDTO) (*res.BaseResponse[GetFileProcessHistoryData], error) {
	fps, pagination, err := s.service.GetFileProcessHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := fromInternalToGetFileHistoryData(fps, pagination)
	resp2 := res.ToResponse(resp)
	return resp2, nil
}

func validateParameterAndSetDataValue(r *http.Request) (*fileprocessing.GetFileProcessHistoryDTO, error) {
	data := &fileprocessing.GetFileProcessHistoryDTO{}
	data.InitDefaultValue()

	values := r.URL.Query()
	for k, v := range values {
		if len(v) > 1 {
			return nil, fmt.Errorf("parameter cannot have multiple value")
		}

		if k == "clientId" {
			data.ClientId = v[0]
		} else {
			val, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("invalid parameter")
			}
			if k == "page" {
				if val == 0 || val > constant.PaginationMaxPage {
					return nil, fmt.Errorf("page out of range")
				}
				data.Page = val
			} else if k == "size" {
				if val == 0 || val > constant.PaginationMaxSize {
					return nil, fmt.Errorf("size out of range")
				}
				data.Size = val
			}
		}
	}

	if data.ClientId == "" {
		return nil, fmt.Errorf("missing clientId")
	}

	return data, nil
}

func (s *Server) CreateProcessByFileAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Bind data & validate input
		data := &CreateFileProcessingRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, commonError.ErrInvalidRequest(err))
			return
		}

		// 2. Handle request
		res, err := s.CreateProcessingFile(r.Context(), data)
		if err != nil {
			render.Render(w, r, commonError.ErrInvalidRequest(err))
			return
		}

		// 3. Render response
		render.Render(w, r, res)
	}
}

func (s *Server) CreateProcessingFile(ctx context.Context, request *CreateFileProcessingRequest) (*res.BaseResponse[CreateFileProcessingResponse], error) {

	// 1. Validate request
	// 2. Call function of Service
	fp, err := s.service.CreateFileProcessing(ctx, &fileprocessing.CreateFileProcessingReqDTO{
		ClientID:    request.ClientID,
		FileURL:     request.FileURL,
		DisplayName: request.FileDisplayName,
		CreatedBy:   request.CreatedBy,
	})
	if err != nil {
		logger.Errorf("CreateProcessingFile: cannot create file processing, got: %v", err)
		return nil, err
	}

	// 3. Map result to response
	resp := &CreateFileProcessingResponse{
		ProcessFileID: int64(fp.ProcessFileID),
	}
	resp2 := res.ToResponse(resp)
	return resp2, nil
}

var _ IServer = &Server{}

// InitFileProcessingServer ...
func InitFileProcessingServer(db *sql.DB) *Server {
	repo := fileprocessing.NewRepo(db)
	service := fileprocessing.NewService(repo)
	return &Server{
		service: service,
	}
}
