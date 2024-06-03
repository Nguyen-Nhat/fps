update config_task set request_body = '[{"field":"sellerId","type":"integer","valuePattern":"$param.sellerId","required":true},{"field":"sellerSku","type":"string","valuePattern":"$A","required":true},{"field":"supplierId","type":"string","valuePattern":"$func.convertSupplierCode2SupplierId;$param.sellerId;$C","required":true},{"field":"siteIds","type":"string","valuePattern":"$func.convertSiteCodes2SiteIds;$param.sellerId;$E","required":false},{"field":"orderDayNote","type":"string","valuePattern":"$F","required":false},{"field":"orderDay","type":"string","valuePattern":"$func.convertOrderDay;$F","required":false}]'
where config_mapping_id = 20 and task_index = 2;

update config_mapping set error_column_index = '$G' where id = 20;
