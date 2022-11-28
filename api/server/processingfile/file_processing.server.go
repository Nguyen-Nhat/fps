package fileprocessing

import (
	"context"
	"database/sql"
	"net/http"

	commonError "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	res "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/go-chi/render"
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

var _ IServer = &Server{}

func (s *Server) GetFileProcessHistoryAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
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

// InitFileProcessingServer ...
func InitFileProcessingServer(db *sql.DB) *Server {
	repo := fileprocessing.NewRepo(db)
	service := fileprocessing.NewService(repo)
	return &Server{
		service: service,
	}
}
