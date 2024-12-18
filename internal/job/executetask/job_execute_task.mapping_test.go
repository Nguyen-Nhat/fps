package executetask

import (
	"context"
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tools/i18n"
)

func Test_getTaskConfig(t *testing.T) {
	configMapping := configloader.ConfigMappingMD{
		Tasks: []configloader.ConfigTaskMD{{TaskIndex: 1}, {TaskIndex: 2}, {TaskIndex: 3}, {TaskIndex: 4, IsAsync: true}},
	}

	type args struct {
		taskIndex     int
		configMapping configloader.ConfigMappingMD
	}
	tests := []struct {
		name        string
		args        args
		want        configloader.ConfigTaskMD
		wantExisted bool
	}{
		{"test getTaskConfig case no config mapping -> not existed", args{1, configloader.ConfigMappingMD{}},
			configloader.ConfigTaskMD{}, false},
		{"test getTaskConfig case task index is 1 -> existed", args{1, configMapping},
			configloader.ConfigTaskMD{TaskIndex: 1}, true},
		{"test getTaskConfig case task index is 2 -> existed", args{2, configMapping},
			configloader.ConfigTaskMD{TaskIndex: 2}, true},
		{"test getTaskConfig case task index is 3 -> existed", args{3, configMapping},
			configloader.ConfigTaskMD{TaskIndex: 3}, true},
		{"test getTaskConfig case task index is 4 -> not existed", args{4, configMapping},
			configloader.ConfigTaskMD{TaskIndex: 4, IsAsync: true}, true},
		{"test getTaskConfig case task index is 5 -> not existed", args{5, configMapping},
			configloader.ConfigTaskMD{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getTaskConfig(tt.args.taskIndex, tt.args.configMapping)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTaskConfig() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantExisted {
				t.Errorf("getTaskConfig() got1 = %v, want %v", got1, tt.wantExisted)
			}
		})
	}
}

func Test_getValueStringFromConfig(t *testing.T) {
	_, _ = i18n.LoadI18n("../../../resources/messages")
	ctx := context.Background()
	previousResponses := map[int32]string{
		1: `{"data":{"name": "this is name", "empty_field": ""}"}`,
		2: `{"data":[{"name": "this is name", "field_int": 12}"}]`,
		3: `{"data":[{"arr": {"field_bool": true}}"}]`,
	}

	type args struct {
		processingFileRow *fileprocessingrow.ProcessingFileRow
		reqField          *configloader.RequestFieldMD
		previousResponses map[int32]string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// ValueDependsOn not support
		{"test getValueStringFromConfig case depend on not support type", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnNone,
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.empty_field",
				Required:             false,
				Type:                 "string",
			}, previousResponses},
			nil, true},
		{"test getValueStringFromConfig case depend on not support type", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnExcel,
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.empty_field",
				Required:             false,
				Type:                 "string",
			}, previousResponses},
			nil, true},
		{"test getValueStringFromConfig case depend on not support type", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnParam,
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.empty_field",
				Required:             false,
				Type:                 "string",
			}, previousResponses},
			nil, true},

		// ValueDependsOn = TASK
		{"test getValueStringFromConfig case depend on existed task, key is not required, existed and empty", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnTask,
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.empty_field",
				Required:             false,
				Type:                 "string",
			}, previousResponses},
			nil, false},
		// case has key
		{"test getValueByPreviousTaskResponse case depend on existed task, key is existed with type int", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnTask,
				ValueDependsOnTaskID: 2,
				ValueDependsOnKey:    "data.0.field_int",
				Required:             true,
				Type:                 "integer",
			}, previousResponses},
			int64(12), false},
		{"test getValueByPreviousTaskResponse case depend on existed task, key is existed with type bool", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnTask,
				ValueDependsOnTaskID: 3,
				ValueDependsOnKey:    "data.0.arr.field_bool",
				Required:             true,
				Type:                 "boolean",
			}, previousResponses},
			true, false},

		// ValueDependsOn = FUNC
		{"test getValueByPreviousTaskResponse case depend on not existed func", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn: configloader.ValueDependsOnFunc,
				ValueDependsOnFunc: customFunc.CustomFunction{
					FunctionPattern: "$func.thisIsNotExistedFunction",
					Name:            "thisIsNotExistedFunction",
					ParamsMapped:    []string{},
				},
				Required: true,
				Type:     "string",
			}, previousResponses},
			nil, true},
		{"test getValueByPreviousTaskResponse case depend on existed func, but func run error", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn: configloader.ValueDependsOnFunc,
				ValueDependsOnFunc: customFunc.CustomFunction{
					FunctionPattern: "$func.testFuncError",
					Name:            "testFuncError",
					ParamsMapped:    []string{},
				},
				Required: true,
				Type:     "string",
			}, previousResponses},
			nil, true},
		{"test getValueByPreviousTaskResponse case depend on existed func, but missing param", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn: configloader.ValueDependsOnFunc,
				ValueDependsOnFunc: customFunc.CustomFunction{
					FunctionPattern: "$func.testFunc;1",
					Name:            "testFunc",
					ParamsMapped:    []string{"1"},
				},
				Required: true,
				Type:     "integer",
			}, previousResponses},
			nil, true},
		{"test getValueByPreviousTaskResponse case depend on existed func, and func run success", args{
			nil,
			&configloader.RequestFieldMD{
				ValueDependsOn:    configloader.ValueDependsOnFunc,
				ValueDependsOnKey: "$func.testFunc;1;5",
				ValueDependsOnFunc: customFunc.CustomFunction{
					FunctionPattern: "$func.testFunc;1;5",
					Name:            "testFunc",
					ParamsMapped:    []string{"1", "5"},
				},
				Required: true,
				Type:     "integer",
			}, previousResponses},
			6, false},
		// ValueDependsOn = DB
		{"test getValueByPreviousTaskResponse case depend on not existed db", args{
			&fileprocessingrow.ProcessingFileRow{
				ProcessingFileRow: ent.ProcessingFileRow{
					ID:     1,
					FileID: 10,
				},
			},
			&configloader.RequestFieldMD{
				ValueDependsOn:    configloader.ValueDependsOnDb,
				ValueDependsOnKey: "testDbField",
				Required:          true,
				Type:              "string",
			},
			nil,
		}, nil, true},
		{"test getValueByPreviousTaskResponse case depend on task id", args{
			&fileprocessingrow.ProcessingFileRow{
				ProcessingFileRow: ent.ProcessingFileRow{
					ID:     1,
					FileID: 10,
				},
			},
			&configloader.RequestFieldMD{
				ValueDependsOn:    configloader.ValueDependsOnDb,
				ValueDependsOnKey: configloader.ValueDependsOnDbFieldTaskId,
				Required:          true,
				Type:              "string",
			},
			nil,
		}, 1, false},
		{"test getValueByPreviousTaskResponse case depend on file id", args{
			&fileprocessingrow.ProcessingFileRow{
				ProcessingFileRow: ent.ProcessingFileRow{
					ID:     1,
					FileID: 10,
				},
			},
			&configloader.RequestFieldMD{
				ValueDependsOn:    configloader.ValueDependsOnDb,
				ValueDependsOnKey: configloader.ValueDependsOnDbFieldFileId,
				Required:          true,
				Type:              "string",
			},
			nil,
		}, int64(10), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValueStringFromConfig(ctx, tt.args.processingFileRow, tt.args.reqField, tt.args.previousResponses, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValueStringFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValueStringFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValueByPreviousTaskResponse(t *testing.T) {
	previousResponses := map[int32]string{
		1: `{"data":{"name": "this is name", "empty_field": ""}"}`,
		2: `{"data":[{"name": "this is name", "field_int": 12}"}]`,
		3: `{"data":[{"arr": {"field_bool": true}}"}]`,
	}

	type args struct {
		reqField          *configloader.RequestFieldMD
		previousResponses map[int32]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test getValueByPreviousTaskResponse case depend on not existed task", args{
			&configloader.RequestFieldMD{
				ValueDependsOnTaskID: 10,
				ValueDependsOnKey:    "data.name",
				Required:             true,
			}, previousResponses},
			"", true},
		{"test getValueByPreviousTaskResponse case task existed, key is required but not existed", args{
			&configloader.RequestFieldMD{
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.nameee",
				Required:             true,
			}, previousResponses},
			"", true},
		{"test getValueByPreviousTaskResponse case task existed, key is not required and not existed", args{
			&configloader.RequestFieldMD{
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.nameee",
				Required:             false,
			}, previousResponses},
			"", true},
		{"test getValueByPreviousTaskResponse case task existed, key is required, existed but empty", args{
			&configloader.RequestFieldMD{
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.empty_field",
				Required:             true,
			}, previousResponses},
			"", true},
		{"test getValueByPreviousTaskResponse case task existed, key is not required, existed and empty", args{
			&configloader.RequestFieldMD{
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.empty_field",
				Required:             false,
			}, previousResponses},
			"", false},
		// case has key
		{"test getValueByPreviousTaskResponse case task existed, key is existed with type int", args{
			&configloader.RequestFieldMD{
				ValueDependsOnTaskID: 2,
				ValueDependsOnKey:    "data.0.field_int",
				Required:             true,
			}, previousResponses},
			"12", false},
		{"test getValueByPreviousTaskResponse case task existed, key is existed with type bool", args{
			&configloader.RequestFieldMD{
				ValueDependsOnTaskID: 3,
				ValueDependsOnKey:    "data.0.arr.field_bool",
				Required:             true,
			}, previousResponses},
			"true", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValueByPreviousTaskResponse(tt.args.reqField, tt.args.previousResponses)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValueByPreviousTaskResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getValueByPreviousTaskResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValueByPreviousTaskResponseForCustomFunc(t *testing.T) {
	type args struct {
		valuePattern      string
		previousResponses map[int32]string
		rowData           []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Test 200 - item string - exist in response", args{
			"$response1.data.products.0.sku",
			map[int32]string{
				1: "{\"code\":0,\"message\":\"Thao tác thành công\",\"data\":{\"products\":[{\"sellerId\":40,\"sku\":\"10125\",\"productLifeCycle\":null,\"isMarketShortage\":null,\"expectedEndOfShortageDate\":null,\"buyable\":null,\"autoReplenishment\":null,\"taxId\":1,\"skus\":[],\"isCoreProductLine\":false}]}}",
			},
			nil,
		}, "10125", false},
		{"Test 200 - item string - exist in response but null", args{
			"$response1.data.products.0.productLifeCycle",
			map[int32]string{
				1: "{\"code\":0,\"message\":\"Thao tác thành công\",\"data\":{\"products\":[{\"sellerId\":40,\"sku\":\"10125\",\"productLifeCycle\":null,\"isMarketShortage\":null,\"expectedEndOfShortageDate\":null,\"buyable\":null,\"autoReplenishment\":null,\"taxId\":1,\"skus\":[],\"isCoreProductLine\":false}]}}",
			},
			nil,
		}, "", false},
		{"Test 400 - item string - not exist in response", args{
			"$response1.data.products.0.test",
			map[int32]string{
				1: "{\"code\":0,\"message\":\"Thao tác thành công\",\"data\":{\"products\":[{\"sellerId\":40,\"sku\":\"10125\",\"productLifeCycle\":null,\"isMarketShortage\":null,\"expectedEndOfShortageDate\":null,\"buyable\":null,\"autoReplenishment\":null,\"taxId\":1,\"skus\":[],\"isCoreProductLine\":false}]}}",
			},
			nil,
		}, "", true},
		{"Test 200 - item string in array - exist in response", args{
			"$response1.data.products.#(sku==\"10125\").sellerId",
			map[int32]string{
				1: "{\"code\":0,\"message\":\"Thao tác thành công\",\"data\":{\"products\":[{\"sellerId\":40,\"sku\":\"10125\",\"productLifeCycle\":null,\"isMarketShortage\":null,\"expectedEndOfShortageDate\":null,\"buyable\":null,\"autoReplenishment\":null,\"taxId\":1,\"skus\":[],\"isCoreProductLine\":false}]}}",
			},
			nil,
		}, "40", false},
		{"Test 200 - item string - exist in response but null", args{
			"$response1.data.products.#(sku==\"10125\").productLifeCycle",
			map[int32]string{
				1: "{\"code\":0,\"message\":\"Thao tác thành công\",\"data\":{\"products\":[{\"sellerId\":40,\"sku\":\"10125\",\"productLifeCycle\":null,\"isMarketShortage\":null,\"expectedEndOfShortageDate\":null,\"buyable\":null,\"autoReplenishment\":null,\"taxId\":1,\"skus\":[],\"isCoreProductLine\":false}]}}",
			},
			nil,
		}, "", false},
		{"Test 400 - item string - not exist in response", args{
			"$response1.data.products.#(sku==\"not_exist\").sellerId",
			map[int32]string{
				1: "{\"code\":0,\"message\":\"Thao tác thành công\",\"data\":{\"products\":[{\"sellerId\":40,\"sku\":\"10125\",\"productLifeCycle\":null,\"isMarketShortage\":null,\"expectedEndOfShortageDate\":null,\"buyable\":null,\"autoReplenishment\":null,\"taxId\":1,\"skus\":[],\"isCoreProductLine\":false}]}}",
			},
			nil,
		}, "", true},
		{"Test 400 - template invalid", args{
			"$response1.",
			map[int32]string{
				1: "{\"code\":0,\"message\":\"Thao tác thành công\",\"data\":{\"products\":[{\"sellerId\":40,\"sku\":\"10125\",\"productLifeCycle\":null,\"isMarketShortage\":null,\"expectedEndOfShortageDate\":null,\"buyable\":null,\"autoReplenishment\":null,\"taxId\":1,\"skus\":[],\"isCoreProductLine\":false}]}}",
			},
			nil,
		}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValueByPreviousTaskResponseForCustomFunc(tt.args.valuePattern, tt.args.previousResponses, tt.args.rowData)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValueByPreviousTaskResponseForCustomFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getValueByPreviousTaskResponseForCustomFunc() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchJsonPath(t *testing.T) {
	type args struct {
		rowData  []string
		jsonPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// case not contains variable
		{"test case path is object and not contains variable", args{[]string{"a", "b", "c"}, `data.transaction.id`},
			`data.transaction.id`},
		{"test case path is array and not contains variable", args{[]string{"a", "b", "c"}, `data.transactions.#(name=="quy")`},
			`data.transactions.#(name=="quy")`},

		// case contains variable
		{"test case path is object and contains variable", args{[]string{"a", "id", "c"}, `data.transaction.{{ $B }}`},
			`data.transaction.id`},
		{"test case path is array and contains variable", args{[]string{"a", "b", "quy"}, `data.transactions.#(name=="{{ $C }}").id`},
			`data.transactions.#(name=="quy").id`},

		// case contains multiple variables
		{"test case path is object and contains multiple variables -> no support all", args{[]string{"a", "id", "member"}, `data.transaction.{{ $C }}.{{ $B }}`},
			`data.transaction.member.{{ $B }}`},
		{"test case path is array and contains multiple variables -> no support all", args{[]string{"id", "b", "quy"}, `data.transactions.#(name=="{{ $C }}").{{ $A }}`},
			`data.transactions.#(name=="quy").{{ $A }}`},

		// case contains space in variable
		{"test case path is array and contains variable", args{[]string{"a", "b", "quy"}, `data.transactions.#(name=="{{               $C                 }}").id`},
			`data.transactions.#(name=="quy").id`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchJsonPath(tt.args.rowData, tt.args.jsonPath)
			if got != tt.want {
				t.Errorf("validateAndMatchJsonPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
