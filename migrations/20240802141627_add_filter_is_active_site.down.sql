-- Client 5
UPDATE config_task t
SET t.task_index   = 4,
    t.request_body = '[{"field":"siteId","type":"integer","valuePattern":"$response1.data.#(sellerSiteCode==\\"{{ $A }}\\").id","required":true},{"field":"binId","type":"integer","valuePattern":"$response3.data.bins.0.binId","required":true},{"field":"actionBy","type":"string","valuePattern":"$param.actionBy","required":true},{"field":"reasonId","type":"integer","valuePattern":"$param.reasonId","required":true},{"field":"items","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$response2.result.products.0.sku","required":true},{"field":"quantity","type":"number","valuePattern":"$E","required":true}]}]'
WHERE t.config_mapping_id = 5 and t.task_index = 3;

UPDATE config_task t
SET t.task_index     = 3,
    t.request_params = '[{"field":"isActive","type":"boolean","valuePattern":"true","required":true},{"field":"siteId","type":"integer","valuePattern":"$response1.data.#(sellerSiteCode==\\"{{ $A }}\\").id","required":true},{"field":"binName","type":"string","valuePattern":"$B","required":true}]'
WHERE t.config_mapping_id = 5 and t.task_index = 2;

UPDATE config_task t
SET t.task_index = 2
WHERE t.config_mapping_id = 5 and t.task_index = 1;

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, method, header, path_params, request_params, request_body, response_success_http_status, response_success_code_schema, response_message_schema, group_by_columns, group_by_size_limit, created_at, updated_at, created_by, message_transformations, is_async) VALUES (5, 1, 'lấy thông tin site', 'http://warehouse3-central-service.warehouse-management/api/v1/sites', 'GET', '', '', '[{"field":"siteIds","type":"integer","valuePattern":"$param.siteId","required":false},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"siteName","type":"string","valuePattern":"$A","required":false}]', '', 200, '{"path":"code","successValues":"0","mustHaveValueInPath":"data.#(sellerSiteCode==\\"{{ $A }}\\").sellerSiteCode"}', '{"path": "message"}', '', 0, '2023-06-23 05:19:00', '2024-04-12 10:04:48', 'quy.tm@teko.vn', '[{"httpStatus":200,"message":""}]', 0);



-- Client 8
UPDATE config_task t
SET t.task_index   = 4,
    t.request_body = '[{"field":"isSubmitOnCreate","type":"boolean","valuePattern":"$param.isNeedSubmit","required":false},{"field":"siteId","type":"integer","valuePattern":"$response1.data.#(sellerSiteCode==\\"{{ $B }}\\").id","required":true},{"field":"note","type":"string","valuePattern":"$D","required":false},{"field":"productStatusType","type":"integer","valuePattern":"$E","required":false},{"field":"contactCode","type":"string","valuePattern":"$F","required":false},{"field":"contactName","type":"string","valuePattern":"$G","required":false},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.actionBy","required":false},{"field":"operationType","type":"integer","valuePattern":"$response2.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).adjustmentTypeId","required":true},{"field":"reasonId","type":"integer","valuePattern":"$response2.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).reasons.#(reasonCode==\\"{{ $C }}\\").reasonId","required":true},{"field":"adjustmentItems","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$response3.result.products.0.sku","required":true},{"field":"quantity","type":"number","valuePattern":"$I","required":true},{"field":"note","type":"string","valuePattern":"$K","required":false},{"field":"imageUrls","type":"array","ArrayItem":[{"field":"","type":"string","valuePattern":"$func.reUploadFile;$J","required":false}]}]}]'
WHERE t.config_mapping_id = 8 and t.task_index = 3;

UPDATE config_task t
SET t.task_index = 3
WHERE t.config_mapping_id = 8 and t.task_index = 2;

UPDATE config_task t
SET t.task_index = 2
WHERE t.config_mapping_id = 8 and t.task_index = 1;

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, method, header, path_params, request_params, request_body, response_success_http_status, response_success_code_schema, response_message_schema, group_by_columns, group_by_size_limit, created_at, updated_at, created_by, message_transformations, is_async) VALUES (8, 1, 'lấy thông tin site', 'http://warehouse3-central-service.warehouse-management/api/v1/sites', 'GET', '', '', '[{"field":"sellerId","type":"string","valuePattern":"$param.sellerId","required":true},{"field":"siteName","type":"string","valuePattern":"$B","required":true},{"field":"siteIds","type":"json","valuePattern":"$param.siteIds","required":false}]', '', 200, '{"path":"code","successValues":"0","mustHaveValueInPath":"data.#(sellerSiteCode==\\"{{ $B }}\\").sellerSiteCode"}', '{"path": "message"}', null, 0, '2024-04-12 09:29:04', '2024-07-10 08:16:13', 'quy.tm@teko.vn', '[{"httpStatus":200,"message":""}]', 0);


-- Client 10
UPDATE config_task t
SET t.task_index   = 3,
    t.request_body = '[{"field":"outboundSiteId","type":"integer","valuePattern":"$response1.data.#(sellerSiteCode==\\"{{ $A }}\\").id","required":true},{"field":"inboundSiteId","type":"integer","valuePattern":"$response2.data.#(sellerSiteCode==\\"{{ $B }}\\").id","required":true},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"createdSource","type":"string","valuePattern":"import_batch_create_screen","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.createdBy","required":true},{"field":"items","type":"json","valuePattern":"$func.convertSellerSkuAndUomName;$C;$param.sellerId","required":true}]'
WHERE config_mapping_id = 10 and task_index = 1;

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, method, header, path_params, request_params, request_body, response_success_http_status, response_success_code_schema, response_message_schema, group_by_columns, group_by_size_limit, created_at, updated_at, created_by, message_transformations, is_async) VALUES (10, 1, 'Lấy thông tin site id từ site code cho kho xuất', 'http://warehouse3-central-service.warehouse-management/api/v1/sites', 'GET', '', '', '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"siteName","type":"string","valuePattern":"$A","required":true}]', '', 200, '{"path":"code","successValues":"0","mustHaveValueInPath":"data.#(sellerSiteCode==\\"{{ $A }}\\").id"}', '{"path": "message"}', '', 0, '2024-06-17 10:33:56', '2024-06-17 10:33:56', 'quy.tm@teko.vn', '[{"httpStatus":0,"message": "Không tìm thấy kho xuất {{$A}}"},{"httpStatus":200,"message":""}]', 0);
INSERT INTO config_task (config_mapping_id, task_index, name, end_point, method, header, path_params, request_params, request_body, response_success_http_status, response_success_code_schema, response_message_schema, group_by_columns, group_by_size_limit, created_at, updated_at, created_by, message_transformations, is_async) VALUES (10, 2, 'Lấy thông tin site id từ site code cho kho nhận', 'http://warehouse3-central-service.warehouse-management/api/v1/sites', 'GET', '', '', '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"siteName","type":"string","valuePattern":"$B","required":true}]', '', 200, '{"path":"code","successValues":"0","mustHaveValueInPath":"data.#(sellerSiteCode==\\"{{ $B }}\\").id"}', '{"path": "message"}', '', 0, '2024-06-17 10:33:56', '2024-06-17 10:33:56', 'quy.tm@teko.vn', '[{"httpStatus":0,"message": "Không tìm thấy kho nhận {{$B}}"},{"httpStatus":200,"message":""}]', 0);


-- Client 17
UPDATE config_task t
SET t.task_index   = 4,
    t.request_body = '[{"field":"isSubmitOnCreate","type":"boolean","valuePattern":"$param.isNeedSubmit","required":false},{"field":"siteId","type":"integer","valuePattern":"$response1.data.#(sellerSiteCode==\\"{{ $B }}\\").id","required":true},{"field":"note","type":"string","valuePattern":"$D","required":false},{"field":"productStatusType","type":"integer","valuePattern":"1","required":false},{"field":"contactCode","type":"string","valuePattern":"","required":false},{"field":"contactName","type":"string","valuePattern":"","required":false},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.actionBy","required":false},{"field":"operationType","type":"integer","valuePattern":"$response2.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).adjustmentTypeId","required":true},{"field":"reasonId","type":"integer","valuePattern":"$response2.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).reasons.#(reasonCode==\\"{{ $C }}\\").reasonId","required":true},{"field":"adjustmentItems","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$response3.result.products.0.sku","required":true},{"field":"quantity","type":"number","valuePattern":"$F","required":true},{"field":"note","type":"string","valuePattern":"$G","required":false},{"field":"unitValue","type":"number","valuePattern":"$I","required":true},{"field":"images","type":"array","ArrayItem":[{"field":"","type":"string","valuePattern":"$func.reUploadFile;$H","required":false}]}]}]'
WHERE t.config_mapping_id = 17 and t.task_index = 3;

UPDATE config_task t
SET t.task_index = 3
WHERE t.config_mapping_id = 17 and t.task_index = 2;

UPDATE config_task t
SET t.task_index = 2
WHERE t.config_mapping_id = 17 and t.task_index = 1;

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, method, header, path_params, request_params, request_body, response_success_http_status, response_success_code_schema, response_message_schema, group_by_columns, group_by_size_limit, created_at, updated_at, created_by, message_transformations, is_async) VALUES (17, 1, 'lấy thông tin site', 'http://warehouse3-central-service.warehouse-management/api/v1/sites', 'GET', '', '', '[{"field":"sellerId","type":"string","valuePattern":"$param.sellerId","required":true},{"field":"siteName","type":"string","valuePattern":"$B","required":true},{"field":"siteIds","type":"json","valuePattern":"$param.siteIds","required":false}]', '', 200, '{"path":"code","successValues":"0","mustHaveValueInPath":"data.#(sellerSiteCode==\\"{{ $B }}\\").sellerSiteCode"}', '{"path": "message"}', '', 0, '2024-03-04 09:42:44', '2024-07-11 03:31:17', 'anh.lt2@teko.vn', '[{"httpStatus":200,"message":""}]', 0);


-- Client 18
UPDATE config_task t
SET t.task_index   = 4,
    t.request_body = '[{"field":"isSubmitOnCreate","type":"boolean","valuePattern":"$param.isNeedSubmit","required":false},{"field":"siteId","type":"integer","valuePattern":"$response1.data.#(sellerSiteCode==\\"{{ $B }}\\").id","required":true},{"field":"note","type":"string","valuePattern":"$D","required":false},{"field":"productStatusType","type":"integer","valuePattern":"9","required":false},{"field":"contactCode","type":"string","valuePattern":"","required":false},{"field":"contactName","type":"string","valuePattern":"","required":false},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.actionBy","required":false},{"field":"operationType","type":"integer","valuePattern":"$response2.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).adjustmentTypeId","required":true},{"field":"reasonId","type":"integer","valuePattern":"$response2.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).reasons.#(reasonCode==\\"{{ $C }}\\").reasonId","required":true},{"field":"adjustmentItems","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$response3.result.products.0.sku","required":true},{"field":"quantity","type":"number","valuePattern":"$F","required":true},{"field":"note","type":"string","valuePattern":"$G","required":false},{"field":"unitValue","type":"number","valuePattern":"$I","required":true},{"field":"images","type":"array","ArrayItem":[{"field":"","type":"string","valuePattern":"$func.reUploadFile;$H","required":false}]}]}]'
WHERE t.config_mapping_id = 18 and t.task_index = 3;

UPDATE config_task t
SET t.task_index = 3
WHERE t.config_mapping_id = 18 and t.task_index = 2;

UPDATE config_task t
SET t.task_index = 2
WHERE t.config_mapping_id = 18 and t.task_index = 1;

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, method, header, path_params, request_params, request_body, response_success_http_status, response_success_code_schema, response_message_schema, group_by_columns, group_by_size_limit, created_at, updated_at, created_by, message_transformations, is_async) VALUES (18, 1, 'lấy thông tin site', 'http://warehouse3-central-service.warehouse-management/api/v1/sites', 'GET', '', '', '[{"field":"sellerId","type":"string","valuePattern":"$param.sellerId","required":true},{"field":"siteName","type":"string","valuePattern":"$B","required":true},{"field":"siteIds","type":"json","valuePattern":"$param.siteIds","required":false}]', '', 200, '{"path":"code","successValues":"0","mustHaveValueInPath":"data.#(sellerSiteCode==\\"{{ $B }}\\").sellerSiteCode"}', '{"path": "message"}', '', 0, '2024-03-04 09:42:44', '2024-07-11 03:31:19', 'anh.lt2@teko.vn', '[{"httpStatus":200,"message":""}]', 0);
