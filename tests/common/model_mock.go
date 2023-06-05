package common

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"time"
)

var processingFiles = []fileprocessing.ProcessingFile{
	{ProcessingFile: ent.ProcessingFile{
		ID:                  1,
		ClientID:            1,
		DisplayName:         "processing_file.xlsx",
		FileURL:             "https://file_incorrect_url.xlsx",
		ResultFileURL:       "https://a.com",
		Status:              fileprocessing.StatusInit,
		FileParameters:      "",
		TotalMapping:        10,
		StatsTotalRow:       100,
		StatsTotalProcessed: 0,
		StatsTotalSuccess:   98,
		ErrorDisplay:        "Success",
		CreatedAt:           time.Now(),
		CreatedBy:           "tan.hm@teko.vn",
		UpdatedAt:           time.Now(),
	}}, {ProcessingFile: ent.ProcessingFile{
		ID:                  2,
		ClientID:            2,
		DisplayName:         "processing_file.xlsx",
		FileURL:             "https://a.com",
		ResultFileURL:       "",
		Status:              fileprocessing.StatusProcessing,
		FileParameters:      "",
		TotalMapping:        2,
		StatsTotalRow:       2,
		StatsTotalProcessed: 0,
		StatsTotalSuccess:   0,
		ErrorDisplay:        "Success",
		CreatedAt:           time.Now(),
		CreatedBy:           "tan.hm@teko.vn",
		UpdatedAt:           time.Now(),
	}}, {ProcessingFile: ent.ProcessingFile{
		ID:                  3,
		ClientID:            1,
		DisplayName:         "processing_file.xlsx",
		FileURL:             "https://a.com",
		ResultFileURL:       "https://a.com",
		Status:              fileprocessing.StatusFailed,
		FileParameters:      "",
		TotalMapping:        10,
		StatsTotalRow:       100,
		StatsTotalProcessed: 0,
		StatsTotalSuccess:   98,
		ErrorDisplay:        "Success",
		CreatedAt:           time.Now(),
		CreatedBy:           "tan.hm@teko.vn",
		UpdatedAt:           time.Now(),
	}}, {ProcessingFile: ent.ProcessingFile{
		ID:                  4,
		ClientID:            10,
		DisplayName:         "processing_file.xlsx",
		FileURL:             "https://a.com",
		ResultFileURL:       "https://a.com",
		Status:              fileprocessing.StatusFinished,
		FileParameters:      "",
		TotalMapping:        10,
		StatsTotalRow:       100,
		StatsTotalProcessed: 0,
		StatsTotalSuccess:   98,
		ErrorDisplay:        "Success",
		CreatedAt:           time.Now(),
		CreatedBy:           "quy.tm@teko.vn",
		UpdatedAt:           time.Now(),
	}}, {ProcessingFile: ent.ProcessingFile{
		ID:                  5,
		ClientID:            10,
		DisplayName:         "processing_file.xlsx",
		FileURL:             "https://file_incorrect_url.xlsx",
		ResultFileURL:       "https://a.com",
		Status:              fileprocessing.StatusInit,
		FileParameters:      "",
		TotalMapping:        10,
		StatsTotalRow:       100,
		StatsTotalProcessed: 0,
		StatsTotalSuccess:   98,
		ErrorDisplay:        "Success",
		CreatedAt:           time.Now(),
		CreatedBy:           "quy.tm@teko.vn",
		UpdatedAt:           time.Now(),
	}}, {ProcessingFile: ent.ProcessingFile{
		ID:                  6,
		ClientID:            12,
		DisplayName:         "processing_file.xlsx",
		FileURL:             "https://file_incorrect_url.xlsx",
		ResultFileURL:       "https://a.com",
		Status:              fileprocessing.StatusInit,
		FileParameters:      "",
		TotalMapping:        10,
		StatsTotalRow:       100,
		StatsTotalProcessed: 0,
		StatsTotalSuccess:   98,
		ErrorDisplay:        "Success",
		CreatedAt:           time.Now(),
		CreatedBy:           "tan.hm@teko.vn",
		UpdatedAt:           time.Now(),
	}}, {ProcessingFile: ent.ProcessingFile{
		ID:                  7,
		ClientID:            12,
		DisplayName:         "processing_file.xlsx",
		FileURL:             "https://file_incorrect_url.xlsx",
		ResultFileURL:       "https://a.com",
		Status:              fileprocessing.StatusInit,
		FileParameters:      "",
		TotalMapping:        10,
		StatsTotalRow:       100,
		StatsTotalProcessed: 0,
		StatsTotalSuccess:   98,
		ErrorDisplay:        "Success",
		CreatedAt:           time.Now(),
		CreatedBy:           "tri.tm1@teko.vn",
		UpdatedAt:           time.Now(),
	}},
}

func GetProcessingFileMockById(id int) *fileprocessing.ProcessingFile {
	for _, file := range processingFiles {
		if file.ID == id {
			return &file
		}
	}
	return nil
}

// Processing File Row -------------------------------------------------------------------------------------------------

var processingFileRows = []fileprocessingrow.ProcessingFileRow{
	{ProcessingFileRow: ent.ProcessingFileRow{
		ID:       1,
		FileID:   2,
		RowIndex: 0,
		RowDataRaw: `
			{"A":"0987747470","B":"Nhà cung cấp 001","C":"seller01@gmaill.com","D":"","E":""}
		`,
		TaskIndex: 1,
		TaskMapping: `
			{"RowId":2,"TaskId":1,"Endpoint":"http://127.0.0.1:8080/api/v4/member/upsert","Header":{"x-client-id":"49795103480352768"},"Request":{"address":{"FieldName":"address","Value":"","IsMappingExcel":true,"IsMappingResponse":false,"MappingKey":"E","DependOnResponseOfTaskId":0},"email":{"FieldName":"email","Value":"","IsMappingExcel":true,"IsMappingResponse":false,"MappingKey":"C","DependOnResponseOfTaskId":0},"idCardNumber":{"FieldName":"idCardNumber","Value":"","IsMappingExcel":true,"IsMappingResponse":false,"MappingKey":"D","DependOnResponseOfTaskId":0},"name":{"FieldName":"name","Value":"","IsMappingExcel":true,"IsMappingResponse":false,"MappingKey":"B","DependOnResponseOfTaskId":0},"phone":{"FieldName":"phone","Value":"","IsMappingExcel":true,"IsMappingResponse":false,"MappingKey":"A","DependOnResponseOfTaskId":0},"requestId":{"FieldName":"requestId","Value":"","IsMappingExcel":true,"IsMappingResponse":false,"MappingKey":"A","DependOnResponseOfTaskId":0}},"Response":{"httpStatusSuccess":"200","code":{"path":"code","successValues":"00"},"message":{"path":"message"}}}
		`,
		TaskDependsOn:  "",
		TaskRequestRaw: "",
		TaskResponseRaw: `
			{"code":"00","message":"Successful","execution_time":57,"result":{"memberQrCode":"a5nm7t2kH6jqTW88YqM97%2Be4XYqlQYwRfrMNNUFFc5yrcwJX62y6%2Bgq95NYOOl01","memberStaticQrCode":"a5nm7t2kH6jqTW88YqM97ws2mXT5QxoSmdO0G6lwJLA%3D","memberId":"380960311550676992","phone":"0987747470","memberCardId":null,"name":"Nhà cung cấp 001","point":57,"tierPoint":0,"isActive":true,"keepTierUntil":1697439616,"accumulationFrom":1609459200,"accumulationTo":1640995200,"currentTierCode":"1_HANGTHUONG","currentTierName":"Hạng Thường","currentTierMinPoint":0,"nextTierCode":"1_HANGBAC","nextTierName":"Hạng Bạc","nextTierMinPoint":10,"defaultCardType":"QR"}}
		`,
		Status:       fileprocessingrow.StatusInit,
		ErrorDisplay: "",
		ExecutedTime: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}}, {ProcessingFileRow: ent.ProcessingFileRow{
		ID:       2,
		FileID:   2,
		RowIndex: 0,
		RowDataRaw: `
			{"A":"0987747470","B":"Nhà cung cấp 001","C":"seller01@gmaill.com","D":"","E":""}
		`,
		TaskIndex: 2,
		TaskMapping: `
			{"RowId":3,"TaskId":2,"Endpoint":"http://127.0.0.1:8080/api/v4/credentials","Header":{"x-client-id":"204477253650747392"},"Request":{"credentialType":{"FieldName":"credentialType","Value":"MEMBER_CARD","IsMappingExcel":false,"IsMappingResponse":false,"MappingKey":"","DependOnResponseOfTaskId":0},"credentialValue":{"FieldName":"credentialValue","Value":"","IsMappingExcel":true,"IsMappingResponse":false,"MappingKey":"A","DependOnResponseOfTaskId":0},"memberId":{"FieldName":"memberId","Value":"","IsMappingExcel":false,"IsMappingResponse":true,"MappingKey":"result.memberId","DependOnResponseOfTaskId":1}},"Response":{"httpStatusSuccess":"","code":{"path":"code","successValues":"00"},"message":{"path":"message"}}}
		`,
		TaskDependsOn:   "",
		TaskRequestRaw:  "",
		TaskResponseRaw: "",
		Status:          fileprocessingrow.StatusInit,
		ErrorDisplay:    "",
		ExecutedTime:    0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}}, {ProcessingFileRow: ent.ProcessingFileRow{ // only for statistic
		ID:              3,
		FileID:          2,
		RowIndex:        1,
		RowDataRaw:      "",
		TaskIndex:       1,
		TaskMapping:     "",
		TaskDependsOn:   "",
		TaskRequestRaw:  "",
		TaskResponseRaw: "",
		Status:          fileprocessingrow.StatusSuccess,
		ErrorDisplay:    "",
		ExecutedTime:    0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}},
	{ProcessingFileRow: ent.ProcessingFileRow{ // only for statistic
		ID:              4,
		FileID:          2,
		RowIndex:        1,
		RowDataRaw:      "",
		TaskIndex:       2,
		TaskMapping:     "",
		TaskDependsOn:   "",
		TaskRequestRaw:  "",
		TaskResponseRaw: "",
		Status:          fileprocessingrow.StatusSuccess,
		ErrorDisplay:    "",
		ExecutedTime:    0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}},
}
