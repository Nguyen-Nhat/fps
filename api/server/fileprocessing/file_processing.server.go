package fileprocessing

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"golang.org/x/text/language"

	commonError "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	res "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tools/i18n"
)

type (
	IServer interface {
		GetFileProcessHistoryAPI() func(http.ResponseWriter, *http.Request)
		CreateProcessByFileAPI() func(http.ResponseWriter, *http.Request)
	}

	// Server ...
	Server struct {
		service fileprocessing.Service
	}
)

// *Server implements IServer
var _ IServer = &Server{}

var supportedLanguagesMatcher = language.NewMatcher([]language.Tag{
	language.Vietnamese, // "vi"
	language.English,    // "en"
})

// InitFileProcessingServer ...
func InitFileProcessingServer(db *sql.DB) *Server {
	repo := fileprocessing.NewRepo(db)
	service := fileprocessing.NewService(repo)
	return &Server{
		service: service,
	}
}

func (s *Server) GetFileProcessHistoryAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.Bind data & validate input
		data := &GetFileProcessHistoryRequest{}
		if err := bindAndValidateRequestParams(r, data); err != nil {
			logger.Errorf("===== GetFileProcessHistoryAPI: Bind data and validate input error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ErrRenderInvalidRequest(err))
			return
		}

		// 2. Handle request
		resp, err := s.GetFileProcessHistory(r.Context(), data)
		if err != nil {
			logger.Errorf("===== GetFileProcessHistoryAPI handler error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}

		// 3. Render response
		err = render.Render(w, r, resp)
		if err != nil {
			logger.Errorf("===== GetFileProcessHistoryAPI render response error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}
	}
}

func (s *Server) GetFileProcessHistory(ctx context.Context, req *GetFileProcessHistoryRequest) (*res.BaseResponse[GetFileProcessHistoryData], error) {
	// 1. Map the request to internal DTO
	input := &fileprocessing.GetFileProcessHistoryDTO{
		ClientID:        req.ClientID,
		ClientIds:       req.ClientIds,
		SellerId:        req.SellerID,
		CreatedBy:       req.CreatedBy,
		Page:            req.Page,
		PageSize:        req.PageSize,
		CreatedByEmails: req.CreatedByEmails,
		ProcessFileIds:  req.ProcessFileIds,
		SearchFileName:  req.SearchFileName,
		MerchantId:      req.MerchantId,
		TenantId:        req.TenantId,
	}

	// 2. Handle request
	fps, pagination, err := s.service.GetFileProcessHistory(ctx, input)
	if err != nil {
		logger.Infof("Error in GetFileProcessHistory Internal")
		return nil, err
	}
	// 3. Return
	resp := fromInternalToGetFileHistoryData(fps, pagination)
	resp2 := res.ToResponse(resp)
	return resp2, nil
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
	if request.SellerID != 0 && request.MerchantId == constant.EmptyString {
		request.MerchantId = strconv.Itoa(int(request.SellerID))
	}

	// 2. Match language
	acceptLanguage := i18n.GetLanguageFromContext(ctx)

	// 3. Call function of Service
	fp, err := s.service.CreateFileProcessing(ctx, &fileprocessing.CreateFileProcessingReqDTO{
		ClientID:       request.ClientID,
		FileURL:        request.FileURL,
		DisplayName:    request.FileDisplayName,
		CreatedBy:      request.CreatedBy,
		FileParameters: request.Parameters,
		SellerID:       request.SellerID,
		MerchantId:     request.MerchantId,
		TenantId:       request.TenantId,
		AcceptLanguage: acceptLanguage,
	})
	if err != nil {
		logger.Errorf("CreateProcessingFile: cannot create file processing, got: %v", err)
		return nil, err
	}

	// 4. Map result to response
	resp := &CreateFileProcessingResponse{
		ProcessFileID: int64(fp.ProcessFileID),
	}
	resp2 := res.ToResponse(resp)
	return resp2, nil
}
