INSERT INTO fps_client (id, client_id, name, description, created_at, created_by, updated_at, sample_file_url)
VALUES (33, 33, 'SC - Import Update Purchase Status', 'SC - Import Update Purchase Status', NOW(), 'anh.lt2@teko.vn',
        NOW(), '');

INSERT INTO config_mapping(id, client_id, total_tasks, data_start_at_row, data_at_sheet, require_column_index,
                           error_column_index,
                           created_at, created_by, updated_at)
VALUES (33, 33, 0, 4, 'Data', '-', '$J', NOW(), 'anh.lt2@teko.vn', NOW());

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, `method`, header, path_params, request_params,
                         request_body,
                         response_success_http_status, response_success_code_schema, response_message_schema,
                         message_transformations, group_by_columns, group_by_size_limit, created_at, created_by,
                         updated_at)
VALUES (33, 1, 'Get Old Purchase Status', 'http://medusa-medusa.supply-chain/api/v1/sku/purchase-status', 'GET', '', '',
        '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"skus","type":"json","valuePattern":"$func.convertSellerSku2Skus;$param.sellerId;$A","required":true}]',
        '',
        200, '{"path":"code","successValues":"0"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":""},{"httpStatus":400,"message":"Lỗi: {{$response.message}}"}]',
        '', 0, NOW(), 'anh.lt2@teko.vn', NOW()),
       (33, 2, 'Get VAT', 'http://medusa-medusa.supply-chain/api/v1/extra-data', 'GET', '', '',
        '[{"field":"groupKeys","type":"string","valuePattern":"vat_import_label","required":true}]',
        '',
        200, '{"path":"code","successValues":"0","mustHaveValueInPath":"data.items.#(key==\\"{{ $H }}\\").label"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":""},{"httpStatus":400,"message":"Lỗi: {{$response.message}}"}]',
        '', 0, NOW(), 'anh.lt2@teko.vn', NOW()),
       (33, 3, 'Update Purchase Status', 'http://medusa-medusa.supply-chain/api/v1/sku/purchase-status', 'POST', '', '',
        '',
        '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"skus","type":"json","valuePattern":"$func.convertSellerSku2Skus;$param.sellerId;$A","required":true},{"field":"productLifeCycle","type":"string","valuePattern":"$func.getValueByPriority;string;$C;$response1.data.products.0.productLifeCycle;normal","required":true},{"field":"buyable","type":"boolean","valuePattern":"$func.getValueByPriority;boolean;$D;$response1.data.products.0.buyable;true","required":true},{"field":"isMarketShortage","type":"boolean","valuePattern":"$func.getValueByPriority;boolean;$E;$response1.data.products.0.isMarketShortage;true","required":true},{"field":"expectedEndOfShortageDate","type":"string","valuePattern":"$func.getValueByPriority;string;$F;$response1.data.products.0.expectedEndOfShortageDate","required":false},{"field":"autoReplenishment","type":"boolean","valuePattern":"$func.getValueByPriority;boolean;$G;$response1.data.products.0.autoReplenishment;true","required":true},{"field":"taxId","type":"integer","valuePattern":"$func.getValueByPriority;integer;$response2.data.items.#(key==\\"{{ $H }}\\").label;$response1.data.products.0.taxId","required":false},{"field":"isCoreProductLine","type":"boolean","valuePattern":"$func.getValueByPriority;boolean;$I;$response1.data.products.0.isCoreProductLine;false","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.createdByEmail","required":true}]',
        200, '{"path":"code","successValues":"0"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":"Thành công"},{"httpStatus":400,"message":"Lỗi: {{$response.message}}"},{"httpStatus":0,"message":"{{$response.message}}"}]',
        '', 0, NOW(), 'anh.lt2@teko.vn', NOW());
