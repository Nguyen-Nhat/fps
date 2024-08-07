-- Client 5
DELETE
FROM config_task
WHERE config_mapping_id = 5 and task_index = 1;

UPDATE config_task t
SET t.task_index = 1
WHERE t.config_mapping_id = 5 and t.task_index = 2;

UPDATE config_task t
SET t.task_index     = 2,
    t.request_params = '[{"field":"isActive","type":"boolean","valuePattern":"true","required":true},{"field":"siteId","type":"integer","valuePattern":"$func.validateAndConvertSiteCode2SiteId;$param.sellerId;$A;$param.siteId","required":true},{"field":"binName","type":"string","valuePattern":"$B","required":true}]'
WHERE t.config_mapping_id = 5 and t.task_index = 3;

UPDATE config_task t
SET t.task_index   = 3,
    t.request_body = '[{"field":"siteId","type":"integer","valuePattern":"$func.validateAndConvertSiteCode2SiteId;$param.sellerId;$A;$param.siteId","required":true},{"field":"binId","type":"integer","valuePattern":"$response2.data.bins.0.binId","required":true},{"field":"actionBy","type":"string","valuePattern":"$param.actionBy","required":true},{"field":"reasonId","type":"integer","valuePattern":"$param.reasonId","required":true},{"field":"items","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$response1.result.products.0.sku","required":true},{"field":"quantity","type":"number","valuePattern":"$E","required":true}]}]'
WHERE t.config_mapping_id = 5 and t.task_index = 4;


-- Client 8
DELETE
FROM config_task
WHERE config_mapping_id = 8 and task_index = 1;

UPDATE config_task t
SET t.task_index = 1
WHERE t.config_mapping_id = 8 and t.task_index = 2;

UPDATE config_task t
SET t.task_index = 2
WHERE t.config_mapping_id = 8 and t.task_index = 3;

UPDATE config_task t
SET t.task_index   = 3,
    t.request_body = '[{"field":"isSubmitOnCreate","type":"boolean","valuePattern":"$param.isNeedSubmit","required":false},{"field":"siteId","type":"integer","valuePattern":"$func.validateAndConvertSiteCode2SiteId;$param.sellerId;$B;$param.siteIds","required":true},{"field":"note","type":"string","valuePattern":"$D","required":false},{"field":"productStatusType","type":"integer","valuePattern":"$E","required":false},{"field":"contactCode","type":"string","valuePattern":"$F","required":false},{"field":"contactName","type":"string","valuePattern":"$G","required":false},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.actionBy","required":false},{"field":"operationType","type":"integer","valuePattern":"$response1.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).adjustmentTypeId","required":true},{"field":"reasonId","type":"integer","valuePattern":"$response1.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).reasons.#(reasonCode==\\"{{ $C }}\\").reasonId","required":true},{"field":"adjustmentItems","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$response2.result.products.0.sku","required":true},{"field":"quantity","type":"number","valuePattern":"$I","required":true},{"field":"note","type":"string","valuePattern":"$K","required":false},{"field":"imageUrls","type":"array","ArrayItem":[{"field":"","type":"string","valuePattern":"$func.reUploadFile;$J","required":false}]}]}]'
WHERE t.config_mapping_id = 8 and t.task_index = 4;

-- Client 10
DELETE
FROM config_task
WHERE config_mapping_id = 10 and task_index = 1;

DELETE
FROM config_task
WHERE config_mapping_id = 10 and task_index = 2;

UPDATE config_task t
SET t.task_index   = 1,
    t.request_body = '[{"field":"outboundSiteId","type":"integer","valuePattern":"$func.convertSiteCode2SiteId;$A;$param.sellerId","required":true},{"field":"inboundSiteId","type":"integer","valuePattern":"$func.convertSiteCode2SiteId;$B;$param.sellerId","required":true},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"createdSource","type":"string","valuePattern":"import_batch_create_screen","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.createdBy","required":true},{"field":"items","type":"json","valuePattern":"$func.convertSellerSkuAndUomName;$C;$param.sellerId","required":true}]'
WHERE config_mapping_id = 10 and task_index = 3;


-- Client 17
DELETE
FROM config_task
WHERE config_mapping_id = 17 and task_index = 1;

UPDATE config_task t
SET t.task_index = 1
WHERE t.config_mapping_id = 17 and t.task_index = 2;

UPDATE config_task t
SET t.task_index = 2
WHERE t.config_mapping_id = 17 and t.task_index = 3;

UPDATE config_task t
SET t.task_index   = 3,
    t.request_body = '[{"field":"isSubmitOnCreate","type":"boolean","valuePattern":"$param.isNeedSubmit","required":false},{"field":"siteId","type":"integer","valuePattern":"$func.validateAndConvertSiteCode2SiteId;$param.sellerId;$B;$param.siteIds","required":true},{"field":"note","type":"string","valuePattern":"$D","required":false},{"field":"productStatusType","type":"integer","valuePattern":"1","required":false},{"field":"contactCode","type":"string","valuePattern":"","required":false},{"field":"contactName","type":"string","valuePattern":"","required":false},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.actionBy","required":false},{"field":"operationType","type":"integer","valuePattern":"$response1.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).adjustmentTypeId","required":true},{"field":"reasonId","type":"integer","valuePattern":"$response1.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).reasons.#(reasonCode==\\"{{ $C }}\\").reasonId","required":true},{"field":"adjustmentItems","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$response2.result.products.0.sku","required":true},{"field":"quantity","type":"number","valuePattern":"$F","required":true},{"field":"note","type":"string","valuePattern":"$G","required":false},{"field":"unitValue","type":"number","valuePattern":"$I","required":true},{"field":"images","type":"array","ArrayItem":[{"field":"","type":"string","valuePattern":"$func.reUploadFile;$H","required":false}]}]}]'
WHERE t.config_mapping_id = 17 and t.task_index = 4;

-- Client 18
DELETE
FROM config_task
WHERE config_mapping_id = 18 and task_index = 1;

UPDATE config_task t
SET t.task_index = 1
WHERE t.config_mapping_id = 18 and t.task_index = 2;

UPDATE config_task t
SET t.task_index = 2
WHERE t.config_mapping_id = 18 and t.task_index = 3;

UPDATE config_task t
SET t.task_index   = 3,
    t.request_body = '[{"field":"isSubmitOnCreate","type":"boolean","valuePattern":"$param.isNeedSubmit","required":false},{"field":"siteId","type":"integer","valuePattern":"$func.validateAndConvertSiteCode2SiteId;$param.sellerId;$B;$param.siteIds","required":true},{"field":"note","type":"string","valuePattern":"$D","required":false},{"field":"productStatusType","type":"integer","valuePattern":"9","required":false},{"field":"contactCode","type":"string","valuePattern":"","required":false},{"field":"contactName","type":"string","valuePattern":"","required":false},{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"createdBy","type":"string","valuePattern":"$param.actionBy","required":false},{"field":"operationType","type":"integer","valuePattern":"$response1.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).adjustmentTypeId","required":true},{"field":"reasonId","type":"integer","valuePattern":"$response1.data.adjustmentTypes.#(reasons.#(reasonCode==\\"{{ $C }}\\")).reasons.#(reasonCode==\\"{{ $C }}\\").reasonId","required":true},{"field":"adjustmentItems","type":"array","ArrayItem":[{"field":"sku","type":"string","valuePattern":"$response2.result.products.0.sku","required":true},{"field":"quantity","type":"number","valuePattern":"$F","required":true},{"field":"note","type":"string","valuePattern":"$G","required":false},{"field":"unitValue","type":"number","valuePattern":"$I","required":true},{"field":"images","type":"array","ArrayItem":[{"field":"","type":"string","valuePattern":"$func.reUploadFile;$H","required":false}]}]}]'
WHERE t.config_mapping_id = 18 and t.task_index = 4;
