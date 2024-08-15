package client_config

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/clientconfig"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fpsclient"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tests/common"
)

type getClientConfigSuite struct {
	suite.Suite
	ctx       context.Context
	db        *sql.DB
	server    *clientconfig.Server
	entClient *ent.Client
	fixedTime time.Time
}

func TestGetClientConfig(t *testing.T) {
	ts := &getClientConfigSuite{}
	ts.ctx = context.Background()
	ts.fixedTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	ts.db, ts.entClient = common.PrepareDatabaseSqlite(ts.ctx, ts.T(), true)
	ts.tearDown()
	ts.server = clientconfig.InitClientConfigServer(ts.db)

	defer func() {
		ts.tearDown()
	}()

	suite.Run(t, ts)
}

func (ts *getClientConfigSuite) tearDown() {
	common.TruncateAllTables(ts.ctx, ts.db, ts.entClient)
}

func (ts *getClientConfigSuite) assert(req *clientconfig.GetClientConfigRequest, wantFile string) {
	defer func() {
		ts.tearDown()
	}()

	g := goldie.New(ts.T())

	resp, err := ts.server.GetClientConfig(ts.ctx, req)
	if err != nil {
		g.Assert(ts.T(), wantFile, []byte(err.Error()))
		return
	}

	assert.Nil(ts.T(), err)
	g.AssertJson(ts.T(), wantFile, resp)
}

func (ts *getClientConfigSuite) Test200_HappyCase_WithUIConfig_ThenReturnSuccess() {
	uiConfig := clientconfig.UIConfig{
		ImportHistoryTable: clientconfig.UIConfigImportHistoryTable{
			IsShowPreviewProcessFile: true,
			IsShowPreviewResultFile:  false,
			IsShowDebug:              true,
			IsShowCreatedBy:          false,
			IsShowReload:             true,
			ColorScheme:              "scheme",
		},
	}
	uiConfigStr, _ := json.Marshal(uiConfig)
	common.MockConfigMapping(ts.ctx, ts.entClient, []configmapping.ConfigMapping{
		{ConfigMapping: ent.ConfigMapping{
			ClientID:              1,
			TenantID:              "OMNI",
			MaxFileSize:           100,
			MerchantAttributeName: "merchant",
			UsingMerchantAttrName: true,
			InputFileType:         "XLSX,CSV",
			UIConfig:              string(uiConfigStr),
			CreatedBy:             "test",
		}},
	})
	common.MockFpsClient(ts.ctx, ts.entClient, []fpsclient.Client{
		{FpsClient: ent.FpsClient{
			ID:                    1,
			ClientID:              1,
			Name:                  "client-test",
			Description:           "description-client-test",
			ImportFileTemplateURL: "http://template",
			CreatedBy:             "test",
		}},
	})

	ts.assert(&clientconfig.GetClientConfigRequest{
		ClientId: 1,
	}, "happy_case_with_ui_config")
}

func (ts *getClientConfigSuite) Test200_HappyCase_WithNilUIConfig_ThenReturnSuccessWithDefault() {
	common.MockConfigMapping(ts.ctx, ts.entClient, []configmapping.ConfigMapping{
		{ConfigMapping: ent.ConfigMapping{
			ClientID:              1,
			TenantID:              "OMNI",
			MaxFileSize:           100,
			MerchantAttributeName: "merchant",
			UsingMerchantAttrName: true,
			InputFileType:         "XLSX,CSV",
			CreatedBy:             "test",
		}},
	})
	common.MockFpsClient(ts.ctx, ts.entClient, []fpsclient.Client{
		{FpsClient: ent.FpsClient{
			ID:                    1,
			ClientID:              1,
			Name:                  "client-test",
			Description:           "description-client-test",
			ImportFileTemplateURL: "http://template",
			CreatedBy:             "test",
		}},
	})

	ts.assert(&clientconfig.GetClientConfigRequest{
		ClientId: 1,
	}, "happy_case_with_nil_ui_config")
}

func (ts *getClientConfigSuite) Test500_ErrorCase_MissingClient_ThenReturnError() {
	common.MockConfigMapping(ts.ctx, ts.entClient, []configmapping.ConfigMapping{
		{ConfigMapping: ent.ConfigMapping{
			ClientID:              1,
			TenantID:              "OMNI",
			MaxFileSize:           100,
			MerchantAttributeName: "merchant",
			UsingMerchantAttrName: true,
			InputFileType:         "XLSX,CSV",
			CreatedBy:             "test",
		}},
	})
	common.MockFpsClient(ts.ctx, ts.entClient, []fpsclient.Client{
		{FpsClient: ent.FpsClient{
			ID:                    1,
			ClientID:              1,
			Name:                  "client-test",
			Description:           "description-client-test",
			ImportFileTemplateURL: "http://template",
			CreatedBy:             "test",
		}},
	})

	ts.assert(&clientconfig.GetClientConfigRequest{}, "error_case_missing_client")
}
