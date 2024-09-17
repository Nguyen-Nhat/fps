-- Client 29 ----
INSERT INTO fps_client (id, client_id, name, description, created_at, created_by, updated_at, sample_file_url)
VALUES (29, 29, 'SC - Import Disable Quotation', 'SC - Import Disable Quotation', NOW(), 'anh.lt2@teko.vn', NOW(), '');

INSERT INTO config_mapping(id, client_id, total_tasks, data_start_at_row, data_at_sheet, require_column_index, error_column_index,
                           created_at, created_by, updated_at)
VALUES (29, 29, 0, 4, 'Data', '-', '$Q', NOW(), 'anh.lt2@teko.vn', NOW());

INSERT INTO config_task (config_mapping_id, task_index, name, end_point, `method`, header, path_params, request_params, request_body,
                         response_success_http_status, response_success_code_schema, response_message_schema,
                         message_transformations, group_by_columns, group_by_size_limit, created_at, created_by,
                         updated_at)
VALUES (29, 1, 'Disable Quotation', 'http://medusa-medusa.supply-chain/api/v1/upsert_supplier_quotation', 'POST', '', '', '',
        '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"supplierId","type":"integer","valuePattern":"$func.convertSupplierCode2SupplierId;$param.sellerId;$A","required":true},{"field":"siteIds","type":"string","valuePattern":"$func.convertSiteCodes2SiteIds;$param.sellerId;$C","required":true},{"field":"isConsignment","type":"boolean","valuePattern":"$func.convertString2Bool;$D","required":true},{"field":"fromDate","type":"string","valuePattern":"$func.convertDateTimeFormat;$E;02/01/2006;2006-01-02T15:04:05Z","required":true},{"field":"toDate","type":"string","valuePattern":"$func.convertDateTimeFormat;$F;02/01/2006;2006-01-02T15:04:05Z","required":true},{"field":"createdSource","type":"string","valuePattern":"import","required":true},{"field":"requestedById","type":"string","valuePattern":"$param.requestedById","required":true},{"field":"action","type":"string","valuePattern":"UPDATE_END_DATE","required":true},{"field":"lineItems","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$func.convertSellerSkuAndUomName2Sku;$param.sellerId;$H;$J","required":true},{"field":"sellerSku","type":"string","valuePattern":"$H","required":true},{"field":"startDate","type":"string","valuePattern":"$func.convertDateTimeFormat;$E;02/01/2006;2006-01-02T15:04:05Z","required":true},{"field":"endDate","type":"string","valuePattern":"$func.convertDateTimeFormat;$F;02/01/2006;2006-01-02T15:04:05Z","required":true}]}]',
        200, '{"path":"code","successValues":"0"}', '{"path": "message"}',
        '[{"httpStatus":200,"message":"Thành công"},{"httpStatus":400,"message":"Lỗi: {{$response.message}}"},{"httpStatus":0,"message":"{{$response.message}}"}]','', 0, NOW(), 'anh.lt2@teko.vn', NOW());
