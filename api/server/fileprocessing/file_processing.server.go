package fileprocessing

import (
	"context"
	"database/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"net/http"

	"fmt"
	commonError "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
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
			render.Render(w, r, commonError.ErrRenderInvalidRequest(err))
			return
		}

		resp, err := s.GetFileProcessHistory(r.Context(), data)
		if err != nil {
			render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}

		logger.Infof("===== Response Get List Upload file: \n%v\n", utils.JsonString(resp))

		render.Render(w, r, resp)
	}
}

func (s *Server) GetFileProcessHistory(ctx context.Context, req *fileprocessing.GetFileProcessHistoryDTO) (*res.BaseResponse[GetFileProcessHistoryData], error) {
	fps, pagination, err := s.service.GetFileProcessHistory(ctx, req)
	if err != nil {
		logger.Infof("Error in GetFileProcessHistory Internal")
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
	logger.Infof("===== Request Get List Upload file: \n%+v\n", values)
	for k, v := range values {
		if len(v) > 1 {
			return nil, fmt.Errorf("parameter cannot have multiple value")
		}

		switch k {
		case "clientId":
			val, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("invalid clientId parameter")
			}
			data.ClientId = int32(val)
		case "page":
			val, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("invalid page parameter")
			}
			if val == 0 || val > constant.PaginationMaxPage {
				return nil, fmt.Errorf("page out of range")
			}
			data.Page = val
		case "size":
			val, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("invalid size parameter")
			}
			if val == 0 || val > constant.PaginationMaxSize {
				return nil, fmt.Errorf("size out of range")
			}
			data.Size = val
		default:

		}
	}

	if data.ClientId == 0 {
		return nil, fmt.Errorf("missing clientId")
	}

	return data, nil
}

func (s *Server) CreateProcessByFileAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Bind data & validate input
		data := &CreateFileProcessingRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, commonError.ErrRenderInvalidRequest(err))
			return
		}

		// 2. Handle request
		res, err := s.CreateProcessingFile(r.Context(), data)
		if err != nil {
			render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}

		logger.Infof("===== Response CreateProcessByFileAPI: \n%v\n", utils.JsonString(res))

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