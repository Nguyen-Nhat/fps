-- Client 31 ----
INSERT INTO fps_client (id, client_id, name, description, created_at, created_by, updated_at, sample_file_url)
VALUES (31, 31, 'Rebate - Import User Config Serial', 'Rebate - Import User Config Serial', NOW(), 'anh.lt2@teko.vn', NOW(), '');

INSERT INTO config_mapping(id, client_id, total_tasks, data_start_at_row, data_at_sheet, require_column_index, error_column_index,
                           created_at, created_by, updated_at)
VALUES (31, 31, 0, 2, 'Data', '-', '$N', NOW(), 'anh.lt2@teko.vn', NOW());

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, `method`, header, path_params, request_params, request_body,
                         response_success_http_status, response_success_code_schema, response_message_schema,
                         message_transformations, group_by_columns, group_by_size_limit, created_at, created_by,
                         updated_at)
VALUES (31, 1, 'Upsert Serial For Rule', 'http://staff-rebate-service-api.bff/api/v1/rule/upsert-serial', 'POST',
        '[{"field":"Authorization","type":"string","valuePattern":"$param.token","required":true}]', '','',
        '[{"field":"ruleId","type":"integer","valuePattern":"$param.ruleId","required":true},{"field":"sellerSku","type":"string","valuePattern":"$A","required":true},{"field":"oldSerial","type":"string","valuePattern":"$C","required":false},{"field":"newSerial","type":"string","valuePattern":"$D","required":false},{"field":"orderCode","type":"string","valuePattern":"$H","required":false},{"field":"supplierCode","type":"string","valuePattern":"$F","required":false},{"field":"billingDate","type":"string","valuePattern":"$I","required":false},{"field":"unitPriceExcludedTax","type":"number","valuePattern":"$K","required":false},{"field":"unitPriceIncludedTax","type":"number","valuePattern":"$M","required":false}]',
        200, '{"path":"code","successValues":"0"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":"Thành công"},{"httpStatus":400,"message":"Lỗi: {{$response.message}}"},{"httpStatus":0,"message":"{{$response.message}}"}]','', 0, NOW(), 'anh.lt2@teko.vn', NOW()
       );
