package clientconfig

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/render"

	commonError "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	res "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/clientconfig"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	IServer interface {
		GetClientConfigAPI() func(http.ResponseWriter, *http.Request)
	}

	// Server ...
	Server struct {
		service clientconfig.Service
	}
)

// *Server implements IServer
var _ IServer = &Server{}

// InitClientConfigServer ...
func InitClientConfigServer(db *sql.DB) *Server {
	return &Server{
		service: clientconfig.NewService(db),
	}
}

func (s *Server) GetClientConfigAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.Bind request params
		data := &GetClientConfigRequest{}
		if err := bindAndValidateRequestParams(r, data); err != nil {
			logger.Errorf("===== GetClientConfigAPI: Bind data and validate input error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ErrRenderInvalidRequest(err))
			return
		}

		// 2. Handle request
		resp, err := s.GetClientConfig(r.Context(), data)
		if err != nil {
			logger.Errorf("===== GetClientConfigAPI handler error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}

		// 3. Render response
		err = render.Render(w, r, resp)
		if err != nil {
			logger.Errorf("===== GetClientConfigAPI render response error: %+v", err.Error())
			_ = render.Render(w, r, commonError.ToErrorResponse(err))
			return
		}
	}
}

func (s *Server) GetClientConfig(ctx context.Context, req *GetClientConfigRequest) (*res.BaseResponse[GetClientConfigResponse], error) {
	// 1. Get fps client config
	clientConfig, err := s.service.GetClientConfigById(ctx, req.ClientId)
	if err != nil {
		logger.Errorf("GetClientConfig: cannot get client config, got: %v", err)
		return nil, err
	}

	// 2. Return
	return res.ToResponse(&GetClientConfigResponse{
		ClientID:              clientConfig.ClientID,
		TenantID:              clientConfig.TenantID,
		MaxFileSize:           clientConfig.MaxFileSize,
		MerchantAttributeName: clientConfig.MerchantAttributeName,
		UsingMerchantAttrName: clientConfig.UsingMerchantAttrName,
		InputFileTypes:        clientConfig.InputFileTypes,
		ImportFileTemplateUrl: clientConfig.ImportFileTemplateUrl,
		UIConfig: UIConfig{
			ImportHistoryTable: UIConfigImportHistoryTable{
				IsShowPreviewProcessFile: clientConfig.UIConfig.ImportHistoryTable.IsShowPreviewProcessFile,
				IsShowPreviewResultFile:  clientConfig.UIConfig.ImportHistoryTable.IsShowPreviewResultFile,
				IsShowDebug:              clientConfig.UIConfig.ImportHistoryTable.IsShowDebug,
				IsShowCreatedBy:          clientConfig.UIConfig.ImportHistoryTable.IsShowCreatedBy,
				IsShowReload:             clientConfig.UIConfig.ImportHistoryTable.IsShowReload,
				ColorScheme:              clientConfig.UIConfig.ImportHistoryTable.ColorScheme,
			},
		},
	}), nil
}
