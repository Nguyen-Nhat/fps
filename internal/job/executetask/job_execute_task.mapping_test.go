package executetask

import (
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
)

func Test_getTaskConfig(t *testing.T) {
	configMapping := configloader.ConfigMappingMD{
		Tasks: []configloader.ConfigTaskMD{{TaskIndex: 1}, {TaskIndex: 2}, {TaskIndex: 3}},
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
		want    interface{}
		wantErr bool
	}{
		// ValueDependsOn not support
		{"test getValueStringFromConfig case depend on not support type", args{
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnNone,
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.empty_field",
				Required:             false,
				Type:                 "string",
			}, previousResponses},
			nil, true},
		{"test getValueStringFromConfig case depend on not support type", args{
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnExcel,
				ValueDependsOnTaskID: 1,
				ValueDependsOnKey:    "data.empty_field",
				Required:             false,
				Type:                 "string",
			}, previousResponses},
			nil, true},
		{"test getValueStringFromConfig case depend on not support type", args{
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
			&configloader.RequestFieldMD{
				ValueDependsOn:       configloader.ValueDependsOnTask,
				ValueDependsOnTaskID: 2,
				ValueDependsOnKey:    "data.0.field_int",
				Required:             true,
				Type:                 "integer",
			}, previousResponses},
			int64(12), false},
		{"test getValueByPreviousTaskResponse case depend on existed task, key is existed with type bool", args{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValueStringFromConfig(tt.args.reqField, tt.args.previousResponses)
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
