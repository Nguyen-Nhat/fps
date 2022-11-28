package fileawardpoint

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	res "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"github.com/go-chi/render"
)

type (
	IServer interface {
		// APIs for Create User ----------------------------------------------------------------------------------------

		// GetDetailAPI is used for defining Router
		GetDetailAPI() func(http.ResponseWriter, *http.Request)
		// GetDetail is called in GetDetailAPI() and it handles logic of API
		GetDetail(context.Context, *GetFileAwardPointDetailRequest) (*res.BaseResponse[GetFileAwardPointDetailResponse], error)

		// APIs for Get List File Award Point --------------------------------------------------------------------------

		GetListAPI() func(http.ResponseWriter, *http.Request)
		GetList(context.Context, *fileawardpoint.GetListFileAwardPointDTO) (*res.BaseResponse[GetListFileAwardPointData], error)

		// APIs for Create List File Award Point --------------------------------------------------------------------------

		CreateFileAwardPointAPI() func(http.ResponseWriter, *http.Request)
		CreateFileAwardPoint(context.Context, *CreateFileAwardPointDetailRequest) (*res.BaseResponse[CreateFileAwardPointDetailResponse], error)

		// APIs for <DO STH> User --------------------------------------------------------------------------------------

		// DoSthUserAPI() ...
		// ...
		// DoSthUserXXX()...
		// ...
	}

	// Server ...
	Server struct {
		service *fileawardpoint.ServiceImpl
	}
)

func (s Server) GetDetailAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Bind data & validate input
		data := &GetFileAwardPointDetailRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, error2.ErrRenderInvalidRequest(err))
			return
		}

		// 2. Handle request
		res, err := s.GetDetail(r.Context(), data)
		if err != nil {
			render.Render(w, r, error2.ToErrorResponse(err))
			return
		}

		// 3. Render response
		render.Render(w, r, res)
	}
}

func (s *Server) GetDetail(ctx context.Context, request *GetFileAwardPointDetailRequest) (*res.BaseResponse[GetFileAwardPointDetailResponse], error) {
	// 1. Map the request to internal DTO if input for Service too complex
	req := &fileawardpoint.GetFileAwardPointDetailDTO{
		Id: request.Id,
	}

	// 2. Call function of Service
	fap, err := s.service.GetFileAwardPoint(ctx, req)
	if err != nil {
		return nil, err
	}

	// 4. Map result to response
	resp := toFapDetailResponseByEntity(fap)
	resp2 := res.ToResponse(resp)
	return resp2, nil
}

func (s *Server) GetListAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := validateAndSetDataValue(r)
		if err != nil {
			render.Render(w, r, error2.ErrRenderInvalidRequest(err))
			return
		}

		resp, err := s.GetList(r.Context(), data)
		if err != nil {
			render.Render(w, r, error2.ToErrorResponse(err))
			return
		}

		render.Render(w, r, resp)
	}
}

func validateAndSetDataValue(r *http.Request) (*fileawardpoint.GetListFileAwardPointDTO, error) {
	data := &fileawardpoint.GetListFileAwardPointDTO{}
	data.InitDefaultValue()

	values := r.URL.Query()
	for k, v := range values {
		if len(v) > 1 {
			return nil, fmt.Errorf("parameter cannot have multiple value")
		}

		val, err := strconv.Atoi(v[0])
		if err != nil {
			return nil, fmt.Errorf("invalid parameter")
		}
		if k == "merchantId" {
			data.MerchantId = val
		} else if k == "page" {
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
	return data, nil
}

func (s *Server) GetList(ctx context.Context, req *fileawardpoint.GetListFileAwardPointDTO) (*res.BaseResponse[GetListFileAwardPointData], error) {
	faps, pagination, err := s.service.GetListFileAwardPoint(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := fromFileArrToGetListData(faps, pagination)
	resp2 := res.ToResponse(resp)
	return resp2, nil
}

func (s *Server) CreateFileAwardPointAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Bind data & validate input
		data := &CreateFileAwardPointDetailRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, error2.ErrRenderInvalidRequest(err))
			return
		}

		// 2. Handle request
		res, err := s.CreateFileAwardPoint(r.Context(), data)
		if err != nil {
			render.Render(w, r, error2.ToErrorResponse(err))
			return
		}

		// 3. Render response
		render.Render(w, r, res)
	}
}

func (s *Server) CreateFileAwardPoint(ctx context.Context, request *CreateFileAwardPointDetailRequest) (*res.BaseResponse[CreateFileAwardPointDetailResponse], error) {

	// 1. Validate request

	fileName := utils.ExtractFileName(request.FileUrl)
	// 2. Call function of Service
	fap, err := s.service.CreateFileAwardPoint(ctx, &fileawardpoint.CreateFileAwardPointReqDTO{
		MerchantID: request.MerchantID,
		FileUrl:    request.FileUrl,
		Note:       request.Note,
		FileName:   fileName.FullName,
	})
	if err != nil {
		logger.Errorf("CreateFileAwardPoint: cannot create file award point, got: %v", err)
		return nil, err
	}

	// 3. Map result to response
	resp := toFapCreateResponseByEntity(fap)
	resp2 := res.ToResponse(resp)
	return resp2, nil
}

var _ IServer = &Server{}

// InitFileAwardPointServer ...
func InitFileAwardPointServer(db *sql.DB) *Server {
	repo := fileawardpoint.NewRepo(db)
	service := fileawardpoint.NewService(repo)
	return &Server{
		service: service,
	}
}
