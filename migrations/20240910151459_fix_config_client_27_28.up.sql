UPDATE config_task
SET request_body = '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"sellerSku","type":"string","valuePattern":"$A","required":true},{"field":"updatedByEmail","type":"string","valuePattern":"$param.createdByEmail","required":true},{"field":"requestedSiteId","type":"integer","valuePattern":"$B","required":true},{"field":"stockRequestStatus","type":"string","valuePattern":"processed","required":true},{"field":"items","type":"json","valuePattern":"$C","required":true}]'
WHERE config_mapping_id = 27 and task_index = 2;

UPDATE config_task
SET request_body = '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"sellerSku","type":"string","valuePattern":"$A","required":true},{"field":"updatedByEmail","type":"string","valuePattern":"$param.createdByEmail","required":true},{"field":"requestedSiteId","type":"integer","valuePattern":"$B","required":true},{"field":"stockRequestStatus","type":"string","valuePattern":"open","required":true},{"field":"items","type":"json","valuePattern":"$C","required":true}]'
WHERE config_mapping_id = 28 and task_index = 2;
