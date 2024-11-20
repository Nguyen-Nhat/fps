-- Client 27 ----
UPDATE config_task SET
   end_point = 'http://staff-service-api.bff/api/v1/submit-stock-requests-to-purchase',
   header = '[{"field":"Authorization","type":"string","valuePattern":"$param.token","required":true}]',
   request_body = '[{"field":"requestedSiteId","type":"integer","valuePattern":"$B","required":true},{"field":"stockRequestStatus","type":"string","valuePattern":"processed","required":true},{"field":"items","type":"json","valuePattern":"$C","required":true}]'
WHERE config_mapping_id = 27 and task_index = 2;
-- Client 28 ----
UPDATE config_task SET
   end_point = 'http://staff-service-api.bff/api/v1/submit-stock-requests-to-purchase',
   header = '[{"field":"Authorization","type":"string","valuePattern":"$param.token","required":true}]',
   request_body = '[{"field":"requestedSiteId","type":"integer","valuePattern":"$B","required":true},{"field":"stockRequestStatus","type":"string","valuePattern":"open","required":true},{"field":"items","type":"json","valuePattern":"$C","required":true}]'
WHERE config_mapping_id = 28 and task_index = 2;