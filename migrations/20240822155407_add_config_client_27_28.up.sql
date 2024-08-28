-- Client 27 ----
INSERT INTO fps_client (id, client_id, name, description, created_at, created_by, updated_at, sample_file_url)
VALUES (27, 27, 'SC - Import Pending PO', 'SC - Import Pending PO', NOW(), 'anh.lt2@teko.vn', NOW(), 'https://docs.google.com/spreadsheets/d/16AoQkTXKznWNN2yDNAkhpvGv9c9EBQfcfh_AXS3rWIg/view');

INSERT INTO config_mapping(id, client_id, total_tasks, data_start_at_row, data_at_sheet, require_column_index, error_column_index,
                           created_at, created_by, updated_at)
VALUES (27, 27, 0, 2, 'data_converted', '-', '$D', NOW(), 'anh.lt2@teko.vn', NOW());

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, `method`, header, path_params, request_params, request_body,
                         response_success_http_status, response_success_code_schema, response_message_schema,
                         message_transformations, group_by_columns, group_by_size_limit, created_at, created_by,
                         updated_at)
VALUES (27, 1, 'Validate sku', 'http://catalog-core-api.catalog/skus', 'GET', '', '',
        '[{"field":"page","type":"integer","valuePattern":"1","required":true},{"field":"pageSize","type":"integer","valuePattern":"1","required":true},{"field":"sellerSkus","type":"string","valuePattern":"$A","required":true},{"field":"sellerIds","type":"string","valuePattern":"$param.sellerId","required":false},{"field":"editingStatusCodes","type":"string","valuePattern":"active,processing","required":false}]',
        '', 200, '{"path": "code", "successValues": "0", "mustHaveValueInPath": "result.products.0.sellerSku"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":""},{"httpStatus":0,"message":"Lỗi: Sku {{$B}} không tồn tại hoặc inactive"}]', '', 0, NOW(), 'anh.lt2@teko.vn', NOW()),
       (27, 2, 'Submit Stock Requests to Purchase', 'http://medusa-medusa.supply-chain/api/v1/submit-stock-requests-to-purchase', 'POST',
        '', '','',
        '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"requestedSiteId","type":"integer","valuePattern":"$B","required":true},{"field":"stockRequestStatus","type":"string","valuePattern":"processed","required":true},{"field":"items","type":"json","valuePattern":"$C","required":true}]',
        200, '{"path":"code","successValues":"0"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":"Thành công"},{"httpStatus":400,"message":"Lỗi: {{$response.message}}"},{"httpStatus":0,"message":"{{$response.message}}"}]','', 0, NOW(), 'anh.lt2@teko.vn', NOW()
       );

-- Client 28 ----
INSERT INTO fps_client (id, client_id, name, description, created_at, created_by, updated_at, sample_file_url)
VALUES (28, 28, 'SC - Import Confirm SR Column V2', 'SC - Import Confirm SR Column V2', NOW(), 'anh.lt2@teko.vn', NOW(), 'https://docs.google.com/spreadsheets/d/1MljJb99caZuMcLgwlhCrqJHzGxynLWiUK5nSKXErgRc/view');

INSERT INTO config_mapping(id, client_id, total_tasks, data_start_at_row, data_at_sheet, require_column_index, error_column_index,
                           created_at, created_by, updated_at)
VALUES (28, 28, 0, 2, 'data_converted', '-', '$D', NOW(), 'anh.lt2@teko.vn', NOW());

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, `method`, header, path_params, request_params, request_body,
                         response_success_http_status, response_success_code_schema, response_message_schema,
                         message_transformations, group_by_columns, group_by_size_limit, created_at, created_by,
                         updated_at)
VALUES (28, 1, 'Validate sku', 'http://catalog-core-api.catalog/skus', 'GET', '', '',
        '[{"field":"page","type":"integer","valuePattern":"1","required":true},{"field":"pageSize","type":"integer","valuePattern":"1","required":true},{"field":"sellerSkus","type":"string","valuePattern":"$A","required":true},{"field":"sellerIds","type":"string","valuePattern":"$param.sellerId","required":false},{"field":"editingStatusCodes","type":"string","valuePattern":"active,processing","required":false}]',
        '', 200, '{"path": "code", "successValues": "0", "mustHaveValueInPath": "result.products.0.sellerSku"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":""},{"httpStatus":0,"message":"Lỗi: Sku {{$B}} không tồn tại hoặc inactive"}]', '', 0, NOW(), 'anh.lt2@teko.vn', NOW()),
       (28, 2, 'Submit Stock Requests to Purchase', 'http://medusa-medusa.supply-chain/api/v1/submit-stock-requests-to-purchase', 'POST',
        '', '','',
        '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"requestedSiteId","type":"integer","valuePattern":"$B","required":true},{"field":"stockRequestStatus","type":"string","valuePattern":"open","required":true},{"field":"items","type":"json","valuePattern":"$C","required":true}]',
        200, '{"path":"code","successValues":"0"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":"Thành công"},{"httpStatus":400,"message":"Lỗi: {{$response.message}}"},{"httpStatus":0,"message":"{{$response.message}}"}]','', 0, NOW(), 'anh.lt2@teko.vn', NOW()
       );
